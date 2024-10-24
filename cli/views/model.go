package views

import (
	"github.com/bardic/gocrib/model"
	"github.com/bardic/gocrib/queries"
	tea "github.com/charmbracelet/bubbletea"
)

type IController interface {
	GetView() IView
}

type IView interface {
	ParseInput(tea.KeyMsg) tea.Msg
	View() string
	Update(msg tea.Msg) tea.Cmd
}

type IHandView interface {
	IView
	BuildHeader() string
	BuildFooter() string
}

type HandModel struct {
	CurrentTurnPlayerId int32
	SelectedCardId      int
	SelectedCardIds     []int32
	Deck                *model.GameDeck
	Player              *queries.Player
	Account             *queries.Account
}
