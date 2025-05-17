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
	Render(*vo.Match, *vo.Deck) string
}

type HandVO struct {
	LocalPlayerID int
	CardIDs       []int
	Deck          *vo.Deck
}

type Tab struct {
	Title    string
	TabState vo.ViewState
}

type ChangeState struct {
	NewState  string
	AccountID int
	MatchID   int
}

type ChangeTab struct {
	TabIndex int
}
