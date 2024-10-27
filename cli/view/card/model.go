package card

import (
	cliVO "cli/vo"
	"queries"
	"vo"
)

type Model struct {
	*cliVO.ViewModel
	*cliVO.HandVO
	State                  queries.Gamestate
	ActiveSlotIndex        int32
	SelectedCardId         int32
	HighlighedId           int32
	HighlightedSlotIndexes []int32
	Deck                   *vo.GameDeck
}
