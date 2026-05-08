package game

import "math/rand"

const CARDS_PER_PLAYER = 2
const NONE = ""

func (s *Session) ShuffleDeck() {
	for i := len(s.Deck) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		s.Deck[i], s.Deck[j] = s.Deck[j], s.Deck[i]
	}
}

func (s *Session) DealInitialCharacters() {
	for _, player := range s.Players {
		player.characters = make([]characters, 0, CARDS_PER_PLAYER)

		for i := 0; i < CARDS_PER_PLAYER; i++ {
			card := s.Deck[0]
			s.Deck = s.Deck[1:]

			player.characters = append(player.characters, card)
		}
	}
}

func (s *Session) NextTurn() {
	s.Turn = (s.Turn + 1) % len(s.PlayerOrder)
	s.PlayedAction = NONE
	s.ActionChallenger = nil
	s.ActionBlocker = nil
	s.TargetPlayer = nil
	s.IsBlockChallenged = false
	s.PassedPlayers = []string{}
	s.FacedUpCard = NONE
	s.InvokeActionAllowed = false
}

func (s *Session) IsDroppedCardEligibleForAction() bool {
	if eligibleAction, ok := characterActions[s.FacedUpCard]; ok {
		return s.PlayedAction == eligibleAction
	}

	return false
}

func (s *Session) IsDroppedCardEligibleToBlock() bool {
	if eligibleBlock, ok := eligibleBlocks[s.FacedUpCard]; ok {
		return s.PlayedAction == eligibleBlock
	}

	return false
}

func (s *Session) FilterActivePlayers() []*Player {
	activePlayers := []*Player{}

	for _, player := range s.Players {
		if len(player.characters) > 0 {
			activePlayers = append(activePlayers, player)
		}
	}

	return activePlayers
}

func (s *Session) PlayerMappedAction() func(*Session, int, []int) {
	actionMap := map[actions]func(*Session, int, []int){
		INCOME:        func(s *Session, _ int, _ []int) { s.CurrentPlayer().Income() },
		FOREIGN_AID:   func(s *Session, _ int, _ []int) { s.CurrentPlayer().ForeignAid() },
		TAX:           func(s *Session, _ int, _ []int) { s.CurrentPlayer().Tax() },
		STEAL:         func(s *Session, _ int, _ []int) { s.CurrentPlayer().Steal(s) },
		EXCHANGE:      func(s *Session, _ int, cardIndices []int) { s.CurrentPlayer().Exchange(s, cardIndices) },
		ASSASSINATION: func(s *Session, cardIdx int, _ []int) { s.CurrentPlayer().Assassination(s, cardIdx) },
		COUP:          func(s *Session, cardIdx int, _ []int) { s.CurrentPlayer().Coup(s, cardIdx) },
	}

	return actionMap[s.PlayedAction]
}

func (s *Session) ActionRequiresTarget() bool {
	switch s.PlayedAction {
		case STEAL, ASSASSINATION, COUP:
			return true
		default:
			return false
	}
}

func (s *Session) ActionRequiresPayment() bool {
	if _, exists := actionPaymentMap[s.PlayedAction]; exists {
		return true
	}

	return false
}

func (s *Session) TargetShouldYieldCard() bool {
	switch s.PlayedAction {
		case ASSASSINATION, COUP:
			return true
		default:
			return false
	}
}

func (s *Session) CanActionBeChallenged() bool {
	switch s.PlayedAction {
		case INCOME, FOREIGN_AID, COUP:
			return false
		default:
			return true
	}
}

func (s *Session) CanActionBeBlocked() bool {
	switch s.PlayedAction {
		case INCOME, TAX, COUP:
			return false
		default:
			return true
	}
}

func (s *Session) InvokeCurrentPlayerAction(selectedCardIndex int, selectedCardIndices []int	) {
	action := s.PlayerMappedAction()
	action(s, selectedCardIndex, selectedCardIndices)
}
