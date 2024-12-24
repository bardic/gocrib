package main

import (
	"github.com/bardic/gocrib/queries/queries"

	"github.com/bardic/gocrib/cli/utils"
	"github.com/bardic/gocrib/cli/view/login"
	cliVO "github.com/bardic/gocrib/cli/vo"
	"github.com/bardic/gocrib/vo"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type CLI struct {
	GameInitd         bool
	ViewStateName     vo.ViewStateName
	ActivePlayerId    int
	account           *queries.Account
	matchId           int32
	currentController cliVO.IController
	isOpponent        bool
	GameMatch         *vo.GameMatch
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
