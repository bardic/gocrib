package card

import (
	"cli/views"
	"model"
	"queries"
)

type CardModel struct {
	*views.ViewModel
	*views.HandModel
	State               queries.Gamestate
	ActiveSlotIdx       int
	SelectedCardId      int
	HighlighedId        int
	HighlightedSlotIdxs []int
	Deck                *model.GameDeck
}
