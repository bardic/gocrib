package vo

import (
	"vo"

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
	Render() string
	ParseInput(tea.KeyMsg) tea.Msg
	Update(msg tea.Msg) tea.Cmd
	GetState() ControllerState
}

type IView interface {
	Render() string
	Init()
}

type IHandView interface {
	IView
	BuildHeader() string
	BuildFooter() string
}

type HandModel struct {
	CurrentTurnPlayerId int32
	CardsToDisplay      []int32
	SelectedCardIds     []int32
	Deck                *vo.GameDeck
}

type IViewModel interface {
}

type ViewModel struct {
	Name string
}

type Controller struct {
	View  IView
	Model IViewModel
}

type Tab struct {
	Title    string
	TabState vo.ViewState
}
