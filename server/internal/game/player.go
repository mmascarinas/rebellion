package game

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/gorilla/websocket"
)

type Player struct {
	ID   				string
	Name 				string
	coins 			int
	characters  []characters
	Conn 				*websocket.Conn
	responseCh  chan string
}

func NewPlayer(id, name string, conn *websocket.Conn) *Player {
	return &Player{
		ID:         id,
		Name:       name,
		Conn:       conn,
		responseCh: make(chan string, 1),
	}
}

func (p *Player) Income() {
	p.coins += 1
}

func (p *Player) ForeignAid() {
	p.coins += 2
}

func (p *Player) Tax() {
	p.coins += 3
}

func (p *Player) Steal(s *Session) {
	target := s.TargetPlayer

	if target.coins == 0 {
		s.Broadcast([]byte(fmt.Sprintf(`{"event":"steal_failed","player":"%s","target":"%s","reason":"no coins"}`,
			p.Name, target.Name)))
		return
	}

	amount := min(target.coins, 2)

	target.coins -= amount
	p.coins += amount

	s.Broadcast([]byte(fmt.Sprintf(`{"event":"steal","player":"%s","target":"%s","amount":%d}`,
		p.Name, target.Name, amount)))
}

func (p *Player) Exchange(s *Session, cardIndices []int) {
	if len(cardIndices) < 1 {
		p.SendPrompt("error", map[string]any{"message": "You must select card(s) to exchange."})
		return
	}

	sort.Slice(cardIndices, func(i, j int) bool { return cardIndices[i] > cardIndices[j] })
	removed := make([]characters, 0, len(cardIndices))

	for _, cardIndex := range cardIndices {
		if cardIndex < 0 || cardIndex >= len(p.characters) {
			p.SendPrompt("error", map[string]any{"message": "Invalid card index."})
			return
		}
		removed = append(removed, p.characters[cardIndex])
		p.characters = append(p.characters[:cardIndex], p.characters[cardIndex+1:]...)
	}

	newCards := make([]characters, 0, len(cardIndices))
	for range cardIndices {
		if len(s.Deck) == 0 {
			p.SendPrompt("error", map[string]any{"message": "Deck is empty; cannot complete exchange."})
			return
		}
		card := s.Deck[0]
		s.Deck = s.Deck[1:]
		newCards = append(newCards, card)
	}

	if len(cardIndices) == 1 {
		idx := cardIndices[0]
		if idx == 0 {
			p.characters = append(newCards, p.characters...)
		} else {
			p.characters = append(p.characters, newCards...)
		}
	} else {
		p.characters = append(p.characters, newCards...)
	}

	s.Deck = append(s.Deck, removed...)
	s.ShuffleDeck()
}

func (p *Player) Assassination(s *Session, cardIndex int) {
	if s.TargetPlayer == nil {
		p.SendPrompt("error", map[string]any{"message": "No target selected for assassination."})
		return
	}

	if cardIndex < 0 || cardIndex >= len(s.TargetPlayer.characters) {
		p.SendPrompt("error", map[string]any{"message": "Invalid card index for assassination."})
		return
	}
	
	removed := s.TargetPlayer.characters[cardIndex]
	s.Deck = append(s.Deck, removed)
	s.TargetPlayer.characters = append(s.TargetPlayer.characters[:cardIndex], s.TargetPlayer.characters[cardIndex+1:]...)
}

func (p *Player) SkipCounterAction(s *Session) {
	s.PassedPlayers = append(s.PassedPlayers, p.ID)
}

func (p *Player) ReplaceShownPlayerCard(s *Session) {
	s.Deck = append(s.Deck, s.FacedUpCard)
	s.FacedUpCard = NONE
	s.ShuffleDeck()

	card := s.Deck[0]
	s.Deck = s.Deck[1:]
	p.characters = append(p.characters, card)
}

func (p *Player) DropCard(s *Session, cardIndex int) {
	if cardIndex < 0 || cardIndex >= len(p.characters) {
		p.SendPrompt("error", map[string]any{"message": "Invalid card index."})
		return
	}

	s.FacedUpCard = p.characters[cardIndex]
	p.characters = append(p.characters[:cardIndex], p.characters[cardIndex+1:]...)
}

func (p *Player) YieldCard(s *Session) {
	s.Deck = append(s.Deck, s.FacedUpCard)
	s.FacedUpCard = NONE
	s.ShuffleDeck()
}

func (p *Player) CanPayAction(s *Session) bool{
	if p.coins >= actionPaymentMap[s.PlayedAction] {
		return true
	}

	return false
}

func (p *Player) Coup(s *Session, cardIndex int) {
	if s.TargetPlayer == nil {
		return
	}

	if cardIndex < 0 || cardIndex >= len(s.TargetPlayer.characters) {
		return
	}

	removed := s.TargetPlayer.characters[cardIndex]

	s.TargetPlayer.characters = append(
		s.TargetPlayer.characters[:cardIndex],
		s.TargetPlayer.characters[cardIndex+1:]...,
	)

	s.Deck = append(s.Deck, removed)

	s.ShuffleDeck()
}

func (p *Player) ChallengeAction(s *Session, actionResponse string) {
	s.ActionChallenger = p
	currentPlayer := s.CurrentPlayer()

	s.Broadcast([]byte(fmt.Sprintf(`{"event":"action_challenged","challenger":"%s","action":"%s","player":"%s"}`,
		s.ActionChallenger.Name, s.PlayedAction, currentPlayer.Name)))

	showCardIndex := PromptPlayer(currentPlayer, "select_card", map[string]any{
		"message": fmt.Sprintf("Select a card index to prove eligibility for %v action", s.PlayedAction),
	}, func(input string) (int, error) {
		return ParseTargetCharacterIndex(input, s, currentPlayer)
	})

	currentPlayer.DropCard(s, showCardIndex)

	if !s.IsDroppedCardEligibleForAction() {
		s.InvokeActionAllowed = false
		s.Broadcast([]byte(fmt.Sprintf(`{"event":"challenge_success","player":"%s","action":"%s"}`,
			currentPlayer.Name, s.PlayedAction)))
		currentPlayer.YieldCard(s)
	}

	if s.IsDroppedCardEligibleForAction() {
		s.InvokeActionAllowed = true
		s.Broadcast([]byte(fmt.Sprintf(`{"event":"challenge_failed","player":"%s","action":"%s"}`,
			currentPlayer.Name, s.PlayedAction)))

		selectedCardIndex := PromptPlayer(s.ActionChallenger, "select_card", map[string]any{
			"message": "Select a card to yield",
		}, func(input string) (int, error) {
			return ParseTargetCharacterIndex(input, s, s.ActionChallenger)
		})

		s.ActionChallenger.DropCard(s, selectedCardIndex)
		s.ActionChallenger.YieldCard(s)
		currentPlayer.ReplaceShownPlayerCard(s)
	}
}

func (p *Player) BlockAction(s *Session) bool {
	s.ActionBlocker = p
	currentPlayer := s.CurrentPlayer()

	s.Broadcast([]byte(fmt.Sprintf(`{"event":"action_blocked","blocker":"%s","action":"%s","player":"%s"}`,
		s.ActionBlocker.Name, s.PlayedAction, currentPlayer.Name)))

	blockCounterResponse := PromptPlayer(currentPlayer, "block_counter_response", map[string]any{
		"message": fmt.Sprintf("%s blocked your %s action. Enter your response", s.ActionBlocker.Name, s.PlayedAction),
	}, func(input string) (string, error) {
		return IsValidBlockCounterResponse(input, s)
	})

	if counterActions(blockCounterResponse) == YIELD {
		s.Broadcast([]byte(fmt.Sprintf(`{"event":"block_accepted","player":"%s","blocker":"%s"}`,
			currentPlayer.Name, s.ActionBlocker.Name)))
		s.InvokeActionAllowed = false
		return true
	}

	if counterActions(blockCounterResponse) == CHALLENGE {
		s.IsBlockChallenged = true
		s.Broadcast([]byte(fmt.Sprintf(`{"event":"block_challenged","challenger":"%s","blocker":"%s"}`,
			currentPlayer.Name, s.ActionBlocker.Name)))

		showCardIndex := PromptPlayer(s.ActionBlocker, "select_card", map[string]any{
			"message": fmt.Sprintf("Select a card index to prove eligibility to block %s action", s.PlayedAction),
		}, func(input string) (int, error) {
			return ParseTargetCharacterIndex(input, s, s.ActionBlocker)
		})

		s.ActionBlocker.DropCard(s, showCardIndex)

		if s.IsDroppedCardEligibleToBlock() {
			s.InvokeActionAllowed = false
			s.Broadcast([]byte(fmt.Sprintf(`{"event":"block_challenge_failed","blocker":"%s","action":"%s"}`,
				s.ActionBlocker.Name, s.PlayedAction)))

			selectedCardIndex := PromptPlayer(currentPlayer, "select_card", map[string]any{
				"message": "Select a card to yield",
			}, func(input string) (int, error) {
				return ParseTargetCharacterIndex(input, s, currentPlayer)
			})

			currentPlayer.DropCard(s, selectedCardIndex)
			currentPlayer.YieldCard(s)
			s.ActionBlocker.ReplaceShownPlayerCard(s)
			return true
		} else {
			s.InvokeActionAllowed = true
			s.Broadcast([]byte(fmt.Sprintf(`{"event":"block_challenge_success","blocker":"%s","action":"%s"}`,
				s.ActionBlocker.Name, s.PlayedAction)))
			s.ActionBlocker.YieldCard(s)
		}
	}

	return false
}

func PromptPlayer[T any](player *Player, event string, data map[string]any, parser func(string) (T, error)) T {
	for {
		player.SendPrompt(event, data)
		input := player.WaitForResponse()

		value, err := parser(strings.TrimSpace(input))
		if err != nil {
			player.SendPrompt("error", map[string]any{"message": err.Error()})
			continue
		}

		return value
	}
}

func (p *Player) SendPrompt(event string, data map[string]any) {
	msg := map[string]any{
		"event": event,
	}
	for k, v := range data {
		msg[k] = v
	}
	raw, _ := json.Marshal(msg)
	p.Conn.WriteMessage(websocket.TextMessage, raw)
}

func (p *Player) WaitForResponse() string {
	return <-p.responseCh
}

func (p *Player) ReceiveResponse(input string) {
	p.responseCh <- input
}
