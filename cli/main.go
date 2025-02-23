package main

import (
	"encoding/json"

	"github.com/bardic/gocrib/queries/queries"

	"github.com/bardic/gocrib/cli/services"
	"github.com/bardic/gocrib/cli/styles"
	"github.com/bardic/gocrib/cli/utils"
	"github.com/bardic/gocrib/cli/view/container"
	"github.com/bardic/gocrib/cli/view/lobby"
	"github.com/bardic/gocrib/cli/view/login"
	cliVO "github.com/bardic/gocrib/cli/vo"
	"github.com/bardic/gocrib/vo"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type CLI struct {
	account           *queries.Account
	currentController cliVO.IUIController
}

func main() {
	utils.Logger, _ = utils.NewLogger()
	defer utils.Logger.Sync()

	p := tea.NewProgram(newCLIModel(login.NewLogin()))

	if _, err := p.Run(); err != nil {
		utils.Logger.Sugar().Error(err)
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
	var cmds []tea.Cmd
	var match *vo.GameMatch

	switch msg := msg.(type) {
	case queries.Account:
		cli.account = &msg
		cmds = append(cmds, func() tea.Msg {
			return vo.StateChangeMsg{
				NewState:  vo.LobbyView,
				AccountId: msg.ID,
			}
		})
	case vo.StateChangeMsg:
		switch msg.NewState {
		case vo.LobbyView:
			cli.currentController = lobby.NewLobby(msg)
		case vo.JoinGameView:
			fallthrough
		case vo.CreateGameView:
			resp := services.GetMatchById(msg.MatchId)
			err := json.Unmarshal(resp.([]byte), &match)
			if err != nil {
				utils.Logger.Sugar().Error(err)
			}

			cli.createMatchView(match)
		}
	}

	cmds = append(cmds, cli.currentController.Update(msg))
	return cli, tea.Batch(cmds...)
}

func (cli *CLI) createMatchView(match *vo.GameMatch) {
	player := &vo.GamePlayer{}

	var gameDeck vo.GameDeck
	resp := services.GetPlayerByForMatchAndAccount(match.ID, cli.account.ID)
	err := json.Unmarshal(resp.([]byte), player)

	if err != nil {
		utils.Logger.Sugar().Error(err)
	}

	resp = services.GetDeckByPlayIdAndMatchId(*player.ID, *match.ID)
	err = json.Unmarshal(resp.([]byte), &gameDeck)

	if err != nil {
		utils.Logger.Sugar().Error(err)
	}

	cli.currentController = container.NewController(match, player, &gameDeck)
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
