package main

import (
	"time"

	"github.com/bardic/gocrib/cli/utils"
	"github.com/bardic/gocrib/cli/views"
	"github.com/bardic/gocrib/model"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
)

type AppModel struct {
	views.ViewModel
	currentView  views.IViewState
	hand         []model.Card
	kitty        []model.Card
	play         []model.Card
	gameState    model.GameState
	timer        timer.Model
	timerStarted bool
}

func (m *AppModel) Init() tea.Cmd {
	return nil
}

func main() {

	utils.Logger, _ = utils.NewLogger()
	defer utils.Logger.Sync() // flushes buffer, if any

	p := tea.NewProgram(newModel())

	if _, err := p.Run(); err != nil {
		utils.Logger.Sugar().Error(err)
	}
}

func newModel() *AppModel {
	m := &AppModel{
		ViewModel: views.ViewModel{
			ActiveSlot: 0,
			// ViewStateName:  views.Login,
			Tabs:           model.GameTabNames,
			LobbyTabs:      model.LobbyTabNames,
			HighlighedId:   0,
			HighlightedIds: []int{},
		},
		currentView: views.LoginView{},
		hand:        []model.Card{},
		gameState:   model.NewGameState,
		timer:       timer.NewWithInterval(time.Hour, time.Second*1),
	}
	return m
}
