package main

import (
	"cli/utils"
	"cli/view/login"
	cliVO "cli/vo"
	"queries"
	"vo"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type CLI struct {
	GameInitd         bool
	ViewStateName     vo.ViewStateName
	ActivePlayerId    int
	account           *queries.Account
	matchId           int
	currentController cliVO.IController
}

func main() {
	utils.Logger, _ = utils.NewLogger()
	defer utils.Logger.Sync() // flushes buffer, if any

	p := tea.NewProgram(newModel())

	if _, err := p.Run(); err != nil {
		utils.Logger.Sugar().Error(err)
	}
}

func (m *CLI) Init() tea.Cmd {
	return textinput.Blink
}

func newModel() *CLI {
	m := &CLI{
		currentController: &login.Controller{},
	}

	m.currentController.Init()

	return m
}
