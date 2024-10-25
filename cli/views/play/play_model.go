package play

import (
	"github.com/bardic/gocrib/cli/views"
	"github.com/bardic/gocrib/model"
	"github.com/bardic/gocrib/queries"
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
