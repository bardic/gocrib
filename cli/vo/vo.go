package vo

import (
	"github.com/bardic/gocrib/cli/utils"
	"github.com/bardic/gocrib/vo"

	tea "github.com/charmbracelet/bubbletea"
)

type ControllerState int

const (
	LoginControllerState ControllerState = iota
	LobbyControllerState
	BoardControllerState
)

type IController interface {
	Render(gameMatch *vo.GameMatch) string
	ParseInput(tea.KeyMsg) tea.Msg
	Update(msg tea.Msg, gameMatch *vo.GameMatch) tea.Cmd
	GetState() ControllerState
	GetModel() IModel
	GetName() string
}

type IView interface {
	Render() string
	Init()
	BuildHeader() string
	BuildFooter() string
}

type IHandView interface {
	IView
	BuildHeader() string
	BuildFooter() string
}

type HandVO struct {
	LocalPlayerID int
	CardIds       []int
	Deck          *vo.GameDeck
}

type IModel interface {
	GetMatch() *vo.GameMatch
	GetPlayer() *vo.GamePlayer
	SetMatch(match *vo.GameMatch)
}

type ViewModel struct {
	Name      string
	AccountId *int
	Gamematch *vo.GameMatch
}

func (m *ViewModel) SetMatch(match *vo.GameMatch) {
	m.Gamematch = match
}

func (m *ViewModel) GetMatch() *vo.GameMatch {
	return m.Gamematch
}

func (m *ViewModel) GetPlayer() *vo.GamePlayer {
	return utils.GetPlayerForAccountId(m.AccountId, m.Gamematch)
}

type Tab struct {
	Title    string
	TabState vo.ViewState
}
