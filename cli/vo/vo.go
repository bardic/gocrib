package vo

import (
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
	ParseInput(tea.KeyMsg) tea.Msg
	Update(msg tea.Msg) tea.Cmd
	GetState() ControllerState
	GetName() string
}

type IUIController interface {
	IController
	Render() string
}

type IGameController interface {
	IController
	Render(*vo.GameMatch, *vo.GameDeck) string
}

type HandVO struct {
	LocalPlayerID int
	CardIds       []int
	Deck          *vo.GameDeck
}

type Tab struct {
	Title    string
	TabState vo.ViewState
}
