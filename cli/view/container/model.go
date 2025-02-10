package container

import (
	"github.com/bardic/gocrib/cli/utils"
	cliVO "github.com/bardic/gocrib/cli/vo"
	"github.com/bardic/gocrib/vo"
)

type Model struct {
	cliVO.ViewModel
	Tabs          []cliVO.Tab
	State         vo.ViewState
	States        []vo.ViewState
	Subcontroller cliVO.IController
	ActiveTab     int
}

func NewModel(match *vo.GameMatch, player *vo.GamePlayer) *Model {
	return &Model{
		Tabs: []cliVO.Tab{
			{
				Title:    "Board",
				TabState: vo.BoardView,
			},
			{
				Title:    "Play",
				TabState: vo.PlayView,
			},
			{
				Title:    "Hand",
				TabState: vo.HandView,
			},
			{
				Title:    "Kitty",
				TabState: vo.KittyView,
			},
		},
		ViewModel: cliVO.ViewModel{
			Name:      "Container",
			AccountId: player.Accountid,
			Gamematch: match,
		},

		ActiveTab: 0,
	}
}

func (m *Model) GetSubcontroller() cliVO.IController {
	return m.Subcontroller
}

func (m *Model) GetMatch() *vo.GameMatch {
	return m.Gamematch
}

func (m *Model) GetPlayer() *vo.GamePlayer {
	return utils.GetPlayerForAccountId(m.AccountId, m.Gamematch)
}
