package container

import (
	cliVO "github.com/bardic/gocrib/cli/vo"
	"github.com/bardic/gocrib/vo"
)

type Model struct {
	cliVO.IModel
	Tabs          []cliVO.Tab
	State         vo.ViewState
	States        []vo.ViewState
	Match         *vo.GameMatch
	Subcontroller cliVO.IController
	LocalPlayer   *vo.GamePlayer
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
		Match:       match,
		LocalPlayer: player,
		ActiveTab:   0,
	}
}

func (m *Model) GetSubcontroller() cliVO.IController {
	return m.Subcontroller
}

func (m *Model) GetMatch() *vo.GameMatch {
	return m.Match
}

func (m *Model) GetPlayer() *vo.GamePlayer {
	return m.LocalPlayer
}
