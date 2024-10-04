package views

import (
	"github.com/bardic/cribbagev2/model"
)

type ViewModel struct {
	LobbyViewState   model.LobbyViewState
	GameViewState    model.GameViewState
	ViewState        model.ViewState
	Tabs             []string
	LobbyTabs        []string
	ActiveTab        int
	ActiveLandingTab int
	ActiveSlot       model.CardSlots
	HighlighedId     int
	HighlightedIds   []int
}
