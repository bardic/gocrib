package views

import (
	"github.com/bardic/gocrib/model"
	"github.com/bardic/gocrib/queries"
	tea "github.com/charmbracelet/bubbletea"
)

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
	currentTurnPlayerId int32
	selectedCardId      int
	selectedCardIds     []int32
	cards               []int32
	deck                *model.GameDeck
	player              *queries.Player
	account             *queries.Account
}
