package card

import (
	cliVO "cli/vo"
	"queries"
	"vo"
)

type Model struct {
	*cliVO.ViewModel
	*cliVO.HandModel
	State                  queries.Gamestate
	ActiveSlotIndex        int
	SelectedCardId         int
	HighlighedId           int
	HighlightedSlotIndexes []int32
	Deck                   *vo.GameDeck
}
