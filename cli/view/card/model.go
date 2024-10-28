package card

import (
	cliVO "cli/vo"
	"queries"
	"vo"
)

type Model struct {
	*cliVO.ViewModel
	*cliVO.HandVO
	State           queries.Gamestate
	ActiveSlotIndex int32
	SelectedCardIds []int32
	Deck            *vo.GameDeck
}
