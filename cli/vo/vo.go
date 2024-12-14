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
	Init()
	Render(gameMatch *vo.GameMatch) string
	ParseInput(tea.KeyMsg) tea.Msg
	Update(msg tea.Msg, gameMatch *vo.GameMatch) tea.Cmd
	GetState() ControllerState
}

type IView interface {
	Render(hand []int32) string
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
	LocalPlayerID int32
	CardIds       []int32
	Deck          *vo.GameDeck
}

type IModel interface {
}

type ViewModel struct {
	Name string
}

type Controller struct {
	View  IView
	Model IModel
}

type Tab struct {
	Title    string
	TabState vo.ViewState
}
