package card

import (
	"github.com/bardic/gocrib/vo"

	cliVO "github.com/bardic/gocrib/cli/vo"
)

type Model struct {
	*cliVO.HandVO
	State           string
	ActiveSlotIndex int
	SelectedCardIDs []int
	Name            string
	LocalPlayer     *vo.Player
	ActivePlayerID  int
	GameMatchID     int
}
