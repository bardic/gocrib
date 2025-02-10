package card

import (
	"github.com/bardic/gocrib/queries/queries"

	cliVO "github.com/bardic/gocrib/cli/vo"
)

type Model struct {
	*cliVO.ViewModel
	*cliVO.HandVO
	State           queries.Gamestate
	ActiveSlotIndex int
	SelectedCardIds []int
}
