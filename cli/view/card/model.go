package card

import (
	"github.com/bardic/gocrib/queries/queries"
	"github.com/bardic/gocrib/vo"

	cliVO "github.com/bardic/gocrib/cli/vo"
)

type Model struct {
	*cliVO.HandVO
	State           queries.Gamestate
	ActiveSlotIndex int
	SelectedCardIds []int
	Name            string
	LocalPlayer     *vo.GamePlayer
	GameMatchId     *int
}
