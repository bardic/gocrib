package kitty

import (
	"cli/views"
	"model"
	"queries"
)

type KittyModel struct {
	views.ViewModel
	State               queries.Gamestate
	Cards               []int32
	ActiveSlotIdx       int
	HighlighedId        int
	HighlightedSlotIdxs []int
	Deck                *model.GameDeck
}
