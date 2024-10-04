package views

import (
	"github.com/bardic/cribbagev2/model"
)

type ViewModel struct {
	LobbyViewState   model.LobbyViewState
	GameViewState    model.GameViewState
	ViewStateName    ViewStateName
	Tabs             []string
	LobbyTabs        []string
	ActiveTab        int
	ActiveLandingTab int
	ActiveSlot       model.CardSlots
	HighlighedId     int
	HighlightedIds   []int
}

type IViewState interface {
	Enter()
	View() string
}

type ViewStateName uint

const (
	Login ViewStateName = iota
	Lobby
	Game
	GameOver
)
