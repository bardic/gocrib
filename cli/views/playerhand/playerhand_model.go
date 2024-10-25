package playerhand

import (
	"cli/views"
	"model"
	"queries"
)

type PlayerHandModel struct {
	views.ViewModel
	State               queries.Gamestate
	MatchId             int32
	PlayerId            int32
	Cards               []int32
	ActiveSlotIdx       int
	HighlighedId        int
	HighlightedSlotIdxs []int
	Deck                *model.GameDeck
}
