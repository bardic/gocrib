package main

import (
	"log"
	"time"

	"github.com/bardic/cribbagev2/cli/views"
	"github.com/bardic/cribbagev2/model"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
)

type appModel struct {
	views.ViewModel
	currentView views.IViewState
	hand        []model.Card
	kitty       []model.Card
	play        []model.Card
	gameState   model.GameState
	timer       timer.Model
}

func (m *appModel) Init() tea.Cmd {
	return m.timer.Init()
}

func main() {
	p := tea.NewProgram(newModel())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func newModel() *appModel {
	m := &appModel{
		ViewModel: views.ViewModel{
			ActiveSlot:     model.CardOne,
			ViewStateName:  views.Login,
			Tabs:           model.GameTabNames,
			LobbyTabs:      model.LobbyTabNames,
			HighlighedId:   0,
			HighlightedIds: []int{},
		},
		currentView: views.LoginView{},
		hand:        []model.Card{},
		gameState:   model.WaitingState,
		timer:       timer.NewWithInterval(time.Hour, time.Second*1),
	}
	return m
}
