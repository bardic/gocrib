package card

import (
	"github.com/bardic/gocrib/queries/queries"

	cliVO "github.com/bardic/gocrib/cli/vo"
	"github.com/bardic/gocrib/vo"
)

type Model struct {
	*cliVO.ViewModel
	*cliVO.HandVO
	State           queries.Gamestate
	ActiveSlotIndex int
	SelectedCardIds []int
	Deck            *vo.GameDeck
}
