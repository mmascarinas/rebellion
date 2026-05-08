package game

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

type (
	characters string
	actions    string
	counterActions  string
)

const (
	COURTESAN  characters = "COURTESAN"
	MERCENARY 	 characters = "MERCENARY"
	BARON 		 characters = "BARON"
	CONSUL     characters = "CONSUL"
	ADMIRAL 	 characters = "ADMIRAL"
)

const (
	INCOME 			  actions = "INCOME"
	FOREIGN_AID   actions = "FOREIGN_AID"
	TAX 					actions = "TAX"
	EXCHANGE 			actions = "EXCHANGE"
	STEAL 	      actions = "STEAL"
	ASSASSINATION actions = "ASSASSINATION"
	COUP				  actions = "COUP"
)

const (
	CHALLENGE counterActions = "CHALLENGE"
	BLOCK     counterActions = "BLOCK"
	PASS 			counterActions = "PASS"
	YIELD     counterActions = "YIELD"
)

var characterActions = map[characters]actions{
	MERCENARY:   ASSASSINATION,
	BARON:      TAX,
	CONSUL: 		EXCHANGE,
	ADMIRAL:    STEAL,
}

var eligibleBlocks = map[characters]actions{
	COURTESAN:  ASSASSINATION,
	ADMIRAL:    STEAL,
	CONSUL: 		STEAL,
	BARON:    	FOREIGN_AID,
}

var actionPaymentMap = map[actions]int{
	ASSASSINATION: 3,
	COUP:          7,
}

type Session struct {
	ID         					string
	MaxPlayers 					int
	Players    					map[string]*Player
	PlayerOrder					[]string
	Turn       					int
	PlayedAction 				actions
	ActionChallenger  	*Player
	ActionBlocker 			*Player
	TargetPlayer 				*Player
	IsBlockChallenged 	bool
	PassedPlayers 			[]string
	FacedUpCard 				characters
	InvokeActionAllowed bool
	Deck 								[]characters
	mu         					sync.RWMutex
}

func NewSession(id string, maxPlayers int) *Session {
	return &Session{
		ID:         id,
		MaxPlayers: maxPlayers,
		Players:    make(map[string]*Player),
		Turn:       0,
	}
}

func (s *Session) AddPlayer(p *Player) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.Players) >= s.MaxPlayers {
		return false
	}

	s.Players[p.ID] = p
	s.PlayerOrder = append(s.PlayerOrder, p.ID)
	return true
}

func (s *Session) RemovePlayer(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.Players, id)
	for i, pid := range s.PlayerOrder {
		if pid == id {
			s.PlayerOrder = append(s.PlayerOrder[:i], s.PlayerOrder[i+1:]...)
			break
		}
	}
}

func (s *Session) Broadcast(msg []byte) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, p := range s.Players {
		p.Conn.WriteMessage(websocket.TextMessage, msg)
	}
}

func (s *Session) PlayerCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.Players)
}

func (s *Session) CurrentPlayer() *Player {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Players[s.PlayerOrder[s.Turn]]
}

func (s *Session) StartGame() {
	s.Deck = []characters{
		COURTESAN, COURTESAN, COURTESAN,
		MERCENARY, MERCENARY, MERCENARY,
		BARON, BARON, BARON,
		CONSUL, CONSUL, CONSUL,
		ADMIRAL, ADMIRAL, ADMIRAL,
	}

	s.ShuffleDeck()
	s.DealInitialCharacters()

	s.Turn = rand.Intn(len(s.PlayerOrder))

	starter := s.CurrentPlayer()
	msg, _ := json.Marshal(map[string]any{
		"event":        "game_started",
		"first_player": starter.Name,
	})
	s.Broadcast(msg)

	for {
		activePlayers := s.FilterActivePlayers()

		if len(activePlayers) == 1 {
			winMsg, _ := json.Marshal(map[string]any{
				"event":  "game_over",
				"winner": activePlayers[0].Name,
			})
			s.Broadcast(winMsg)
			break
		}

		currentPlayer := s.CurrentPlayer()

		if len(currentPlayer.characters) == 0 {
			s.NextTurn()
			continue
		}

		turnMsg, _ := json.Marshal(map[string]any{
			"event":  "turn_started",
			"player": currentPlayer.Name,
		})
		s.Broadcast(turnMsg)

		var action string
		var selectedCardIndex int
		var selectedCardIndices []int

		if currentPlayer.coins > 9 {
			action = string(COUP)
			s.Broadcast([]byte(fmt.Sprintf(`{"event":"forced_coup","player":"%s"}`, currentPlayer.Name)))
		} else {
			action = PromptPlayer(currentPlayer, "select_action", map[string]any{
				"message": "Enter your action",
			}, func(input string) (string, error) {
				return IsValidPlayerAction(actions(input))
			})
		}

		s.PlayedAction = actions(action)

		if s.ActionRequiresPayment() {
			if !currentPlayer.CanPayAction(s) {
				currentPlayer.SendPrompt("error", map[string]any{
					"message": fmt.Sprintf("Not enough coins to perform %s", action),
				})
				s.PlayedAction = ""
				continue
			}
			currentPlayer.coins -= actionPaymentMap[s.PlayedAction]
		}

		if s.ActionRequiresTarget() {
			targetIndex := PromptPlayer(currentPlayer, "select_target", map[string]any{
				"message": "Enter the index of your target player",
			}, func(input string) (int, error) {
				return IsValidPlayerTarget(input, s)
			})

			s.TargetPlayer = s.Players[s.PlayerOrder[targetIndex]]
		}

		if s.PlayedAction == COUP {
			selectedCardIndex = PromptPlayer(s.TargetPlayer, "select_card", map[string]any{
				"message": "COUP has been launched against you. Select a card to yield",
			}, func(input string) (int, error) {
				return ParseTargetCharacterIndex(input, s, s.TargetPlayer)
			})

			s.InvokeActionAllowed = true
		}

		if s.PlayedAction == INCOME {
			currentPlayer.Income()
			s.Broadcast([]byte(fmt.Sprintf(`{"event":"action_performed","player":"%s","action":"INCOME","coins":%d}`,
				currentPlayer.Name, currentPlayer.coins)))
			s.InvokeActionAllowed = true
		}

		if s.CanActionBeBlocked() || s.CanActionBeChallenged() {
			for i := (s.Turn + 1) % len(s.PlayerOrder); i != s.Turn; i = (i + 1) % len(s.PlayerOrder) {
				responder := s.Players[s.PlayerOrder[i]]

				if len(responder.characters) == 0 {
					continue
				}

				response := PromptPlayer(responder, "counter_action", map[string]any{
					"message": fmt.Sprintf("Enter your counter action to %s's %s", currentPlayer.Name, action),
				}, func(input string) (string, error) {
					return ValidateCounterAction(input, s)
				})

				if counterActions(response) == PASS {
					responder.SkipCounterAction(s)
					s.Broadcast([]byte(fmt.Sprintf(`{"event":"player_passed","player":"%s"}`, responder.Name)))
					continue
				}

				if counterActions(response) == CHALLENGE {
					responder.ChallengeAction(s, response)
					break
				}

				if counterActions(response) == BLOCK {
					if responder.BlockAction(s) {
						break
					}
				}
			}
		}

		if len(s.PassedPlayers) == len(activePlayers)-1 || s.InvokeActionAllowed {
			if s.TargetShouldYieldCard() && s.PlayedAction != COUP {
				selectedCardIndex = PromptPlayer(s.TargetPlayer, "select_card", map[string]any{
					"message": fmt.Sprintf("You are the target of %s. Select a card to yield", s.PlayedAction),
				}, func(input string) (int, error) {
					return ParseTargetCharacterIndex(input, s, s.TargetPlayer)
				})
			}

			if s.PlayedAction == EXCHANGE {
				selectedCardIndices = PromptPlayer(currentPlayer, "select_exchange_cards", map[string]any{
					"message": "Select cards to replace from the deck (comma-separated indices)",
				}, func(input string) ([]int, error) {
					return ParseSelectedIndices(input, s)
				})
			}

			s.InvokeCurrentPlayerAction(selectedCardIndex, selectedCardIndices)

			s.Broadcast([]byte(fmt.Sprintf(`{"event":"action_resolved","player":"%s","action":"%s"}`,
				currentPlayer.Name, s.PlayedAction)))
		}

		s.NextTurn()
	}
}

func ValidateCounterAction(action string, s *Session) (string, error) {
	if !IsValidCounterAction(action) {
		return "", fmt.Errorf("Valid counter actions are: %v, %v, %v", 
			CHALLENGE, BLOCK, PASS,
		)
	}

	if !s.CanActionBeChallenged() && counterActions(action) == CHALLENGE {
		return "", fmt.Errorf("%v action cannot be challenged. Try Again", s.PlayedAction)
	}
	
	if !s.CanActionBeBlocked() && counterActions(action) == BLOCK {
		return "", fmt.Errorf("%v action cannot be blocked. Try Again", s.PlayedAction)
	}

	return action, nil
}

func ParseTargetCharacterIndex(input string, s *Session, p *Player) (int, error) {
	idx, err := strconv.Atoi(strings.TrimSpace(input))

	if err != nil {
		return -1, fmt.Errorf("Card index value must be an integer: %v", input)
	}

	if idx < 0 || idx >= len(p.characters) {
		return -1, fmt.Errorf("Selected index out of range of %v's characters: %d", p.Name, idx)
	}

	return idx, nil
}

func IsValidPlayerTarget(input string, s *Session) (int, error) {
	targetIndex, err := strconv.Atoi(strings.TrimSpace(input))

	if err != nil {
		return -1, fmt.Errorf("Target player index must be an integer: %v", input)
	}

	if targetIndex < 0 || targetIndex >= len(s.PlayerOrder) {
		return -1, fmt.Errorf("Selected index is out of range of game players: %d", targetIndex)
	}

	target := s.Players[s.PlayerOrder[targetIndex]]

	if target == s.CurrentPlayer() {
		return -1, fmt.Errorf("Target player cannot be yourself: %d", targetIndex)
	}

	if len(target.characters) == 0 {
		return -1, fmt.Errorf("Target player is no longer active: %d", targetIndex)
	}

	return targetIndex, nil
}

func IsValidBlockCounterResponse(input string, s *Session) (string, error) {
	switch counterActions(input) {
		case YIELD, CHALLENGE:
			return input, nil
		default:
			return "", fmt.Errorf("Valid block responses are: %v, %v",
				YIELD, CHALLENGE,
			)
	}
}

func ParseSelectedIndices(input string, s *Session) ([]int, error) {
	selectedCardIndices := []int{}

	parts := strings.Split(input, ",")

	seen := make(map[int]bool)

	for _, p := range parts {
		p = strings.TrimSpace(p)
		idx, err := strconv.Atoi(p)

		if err != nil {
			return nil, fmt.Errorf("Card indeces must all be integers: %v", p)
		}

		if idx < 0 || idx >= len(s.CurrentPlayer().characters) {
			return nil, fmt.Errorf("Selected card indeces out range of %v's characters: %d", s.CurrentPlayer().Name, idx)
		}

		if seen[idx] {
			return nil, fmt.Errorf("Provided card indeces contain duplicates: %d", idx)
		}

		seen[idx] = true
		selectedCardIndices = append(selectedCardIndices, idx)
	}

	return selectedCardIndices, nil
}

func IsValidPlayerAction(action actions) (string, error) {
	switch action {
		case INCOME, FOREIGN_AID, TAX, EXCHANGE, STEAL, ASSASSINATION, COUP:
			return string(action), nil
		default:
			return "", fmt.Errorf("Valid actions are: %v, %v, %v, %v, %v, %v, %v", 
			INCOME, FOREIGN_AID, TAX, EXCHANGE, STEAL, ASSASSINATION, COUP,
		)
	}
}

func IsValidCounterAction(response string) bool {
	switch counterActions(response) {
		case CHALLENGE, BLOCK, PASS:
			return true
		default:
			return false
	}
}
