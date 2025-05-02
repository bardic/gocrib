package main

import (
	"github.com/bardic/gocrib/cli/styles"
	logger "github.com/bardic/gocrib/cli/utils/log"
	"github.com/bardic/gocrib/cli/view/container"
	"github.com/bardic/gocrib/cli/view/lobby"
	"github.com/bardic/gocrib/cli/view/login"
	cliVO "github.com/bardic/gocrib/cli/vo"
	"github.com/bardic/gocrib/vo"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type CLI struct {
	currentController cliVO.IUIController
}

func main() {
	l := logger.Get()
	defer l.Sync()

	p := tea.NewProgram(
		newCLIModel(
			login.NewLogin(),
		),
	)

	if _, err := p.Run(); err != nil {
		l.Sugar().Error(err)
	}
}

func newCLIModel(activeController cliVO.IUIController) *CLI {
	return &CLI{
		currentController: activeController,
	}
}

func (cli *CLI) Init() tea.Cmd {
	return textinput.Blink
}

func (cli *CLI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	l := logger.Get()
	defer l.Sync()

	var cmds []tea.Cmd

	stateMsg := msg.(vo.StateChangeMsg)

	switch stateMsg.NewState {
	case vo.LobbyView:
		cli.currentController = lobby.NewLobby(stateMsg)
	case vo.JoinGameView:
		fallthrough
	case vo.CreateGameView:
		cli.currentController = container.NewController(stateMsg)
	}

	cmds = append(cmds, cli.currentController.Update(stateMsg))
	return cli, tea.Batch(cmds...)
}

func (m *CLI) View() string {
	switch m.currentController.(type) {
	case *login.Controller:
		return styles.ViewStyle.Render(m.currentController.Render())
	case *lobby.Controller:
		return styles.ViewStyle.Render(m.currentController.Render())
	case *container.Controller:
		return styles.ViewStyle.Render(m.currentController.Render())
	default:
		return "No view"
	}
}
