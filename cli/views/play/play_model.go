package play

import (
	"cli/views"
	"model"
	"queries"
)

type PlayModel struct {
	views.ViewModel
	State               queries.Gamestate
	ActiveSlotIdx       int
	HighlighedId        int
	HighlightedSlotIdxs []int
	Cards               []int32
	Deck                *model.GameDeck
}
