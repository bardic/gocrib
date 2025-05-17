package container

import (
	"github.com/bardic/gocrib/cli/utils"
	cliVO "github.com/bardic/gocrib/cli/vo"
	"github.com/bardic/gocrib/vo"
)

type Model struct {
	AccountID     int
	Gamematch     *vo.Match
	GameDeck      *vo.Deck
	Tabs          []cliVO.Tab
	State         vo.ViewState
	States        []vo.ViewState
	Subcontroller cliVO.IGameController
	ActiveTab     int
}

func NewModel(match *vo.Match, player *vo.Player, gameDeck *vo.Deck) *Model {
	return &Model{
		GameDeck:  gameDeck,
		AccountID: player.Accountid,
		Gamematch: match,
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
		ActiveTab: 0,
	}
}

func (m *Model) GetSubcontroller() cliVO.IController {
	return m.Subcontroller
}

func (m *Model) GetMatch() *vo.Match {
	return m.Gamematch
}

func (m *Model) GetPlayer() *vo.Player {
	return utils.GetPlayerForAccountID(m.AccountID, m.Gamematch)
}
