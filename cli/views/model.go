package views

import (
	"github.com/bardic/gocrib/model"
	tea "github.com/charmbracelet/bubbletea"
)

type ViewModel struct {
	LobbyViewState model.LobbyViewState
	GameViewState  model.GameViewState
	// ViewStateName    ViewStateName
	Tabs             []string
	LobbyTabs        []string
	ActiveTab        int
	ActiveLandingTab int
	ActiveSlot       model.CardSlots
	HighlighedId     int
	HighlightedIds   []int
}

type IViewState interface {
	Enter() tea.Msg
	View() string
}

type ViewStateName uint

const (
	Login ViewStateName = iota
	Lobby
	Game
	GameOver
)
