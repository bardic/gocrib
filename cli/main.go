package main

import (
	"time"

	"github.com/bardic/gocrib/cli/utils"
	"github.com/bardic/gocrib/cli/views"
	"github.com/bardic/gocrib/model"
	"github.com/bardic/gocrib/queries"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {

	utils.Logger, _ = utils.NewLogger()
	defer utils.Logger.Sync() // flushes buffer, if any

	p := tea.NewProgram(newModel())

	if _, err := p.Run(); err != nil {
		utils.Logger.Sugar().Error(err)
	}
}

type AppModel struct {
	GameInitd      bool
	ViewStateName  model.ViewStateName
	ActivePlayerId int
	account        *queries.Account
	matchId        int
	currentView    views.IView
	timer          timer.Model
	playersReady   bool
}

func (m *AppModel) Init() tea.Cmd {
	return textinput.Blink
}

func newModel() *AppModel {
	m := &AppModel{
		currentView: &views.LoginView{},
		timer:       timer.NewWithInterval(time.Hour, time.Second*1),
	}

	m.currentView.(*views.LoginView).Init()

	return m
}
