package kitty

import (
	"github.com/bardic/gocrib/cli/views"
	"github.com/bardic/gocrib/queries"
)

type KittyModel struct {
	views.ViewModel
	State         queries.Gamestate
	ActiveSlotIdx int
	Cards         []int32
	HighlighedId  int
}
