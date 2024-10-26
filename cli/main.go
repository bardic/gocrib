package main

import (
	"time"

	"cli/utils"
	"cli/view/login"
	"queries"
	"vo"
	cliVO "cli/vo"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
)

type AppModel struct {
	GameInitd         bool
	ViewStateName     vo.ViewStateName
	ActivePlayerId    int
	account           *queries.Account
	matchId           int
	currentController cliVO.IController
	timer             timer.Model
}

func main() {
	utils.Logger, _ = utils.NewLogger()
	defer utils.Logger.Sync() // flushes buffer, if any

	p := tea.NewProgram(newModel())

	if _, err := p.Run(); err != nil {
		utils.Logger.Sugar().Error(err)
	}
}

func (m *AppModel) Init() tea.Cmd {
	return textinput.Blink
}

func newModel() *AppModel {
	m := &AppModel{
		currentController: &login.LoginController{},
		timer:             timer.NewWithInterval(time.Hour, time.Second*1),
	}

	m.currentController.Init()

	return m
}
