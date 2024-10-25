package kitty

import (
	"github.com/bardic/gocrib/cli/views"
	"github.com/bardic/gocrib/model"
	"github.com/bardic/gocrib/queries"
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
