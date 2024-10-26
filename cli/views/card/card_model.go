package card

import (
	"cli/views"
	"model"
	"queries"
)

type CardModel struct {
	*views.ViewModel
	State               queries.Gamestate
	Cards               []int32
	ActiveSlotIdx       int
	HighlighedId        int
	HighlightedSlotIdxs []int
	Deck                *model.GameDeck
}
