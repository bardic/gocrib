package play

import (
	"github.com/bardic/gocrib/cli/views"
	"github.com/bardic/gocrib/queries"
)

type PlayModel struct {
	views.ViewModel
	State        queries.Gamestate
	ActiveSlotId int
	HighlighedId int
	Cards        []int32
	Deck         *queries.Deck
}
