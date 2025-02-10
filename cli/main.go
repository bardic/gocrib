package main

import (
	"encoding/json"

	"github.com/bardic/gocrib/queries/queries"

	"github.com/bardic/gocrib/cli/services"
	"github.com/bardic/gocrib/cli/styles"
	"github.com/bardic/gocrib/cli/utils"
	"github.com/bardic/gocrib/cli/view/container"
	"github.com/bardic/gocrib/cli/view/game"
	"github.com/bardic/gocrib/cli/view/lobby"
	"github.com/bardic/gocrib/cli/view/login"
	cliVO "github.com/bardic/gocrib/cli/vo"
	"github.com/bardic/gocrib/vo"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type CLI struct {
	account           *queries.Account
	currentController cliVO.IController
}

func main() {
	utils.Logger, _ = utils.NewLogger()
	defer utils.Logger.Sync()

	p := tea.NewProgram(newCLIModel(login.New()))

	if _, err := p.Run(); err != nil {
		utils.Logger.Sugar().Error(err)
	}
}

func newCLIModel(activeController cliVO.IController) *CLI {
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
			cli.currentController = &lobby.Controller{
				Controller: &game.Controller{
					Model: lobby.Model{
						ViewModel: cliVO.ViewModel{
							Name:      "Lobby",
							AccountId: msg.AccountId,
							Gamematch: &vo.GameMatch{},
						},
					},
					View: &lobby.View{},
				},
			}

			services.GetOpenMatches()
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

	cmds = append(cmds, cli.currentController.Update(msg, match))
	return cli, tea.Batch(cmds...)
}

func (cli *CLI) createMatchView(match *vo.GameMatch) {
	player := &vo.GamePlayer{}
	resp := services.GetPlayerByForMatchAndAccount(match.ID, cli.account.ID)
	err := json.Unmarshal(resp.([]byte), player)

	if err != nil {
		utils.Logger.Sugar().Error(err)
	}

	cli.currentController = container.NewController(match, player)
}

func (m *CLI) View() string {
	switch m.currentController.(type) {
	case *login.Controller:
		return styles.ViewStyle.Render(m.currentController.Render(nil))
	case *lobby.Controller:
		return styles.ViewStyle.Render(m.currentController.Render(nil))
	case *container.Controller:
		model := m.currentController.GetModel().(*container.Model)
		match := model.ViewModel.Gamematch
		return styles.ViewStyle.Render(m.currentController.Render(match))
	default:
		return "No view"
	}
}
