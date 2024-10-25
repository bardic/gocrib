package hand

import (
	"github.com/bardic/gocrib/cli/views"
	"github.com/bardic/gocrib/queries"
)

type HandModel struct {
	views.ViewModel
	State               queries.Gamestate
	MatchId             int32
	PlayerId            int32
	Cards               []int32
	ActiveSlotIdx       int
	HighlighedId        int
	HighlightedSlotIdxs []int
}
