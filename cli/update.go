package main

import (
	"encoding/json"
	"fmt"

	"github.com/bardic/gocrib/queries/queries"

	"github.com/bardic/gocrib/cli/services"
	"github.com/bardic/gocrib/cli/utils"
	"github.com/bardic/gocrib/cli/view/container"
	"github.com/bardic/gocrib/cli/view/gameContainer"
	"github.com/bardic/gocrib/cli/view/lobby"
	cliVO "github.com/bardic/gocrib/cli/vo"
	"github.com/bardic/gocrib/vo"

	tea "github.com/charmbracelet/bubbletea"
)

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
		fmt.Println(msg.NewState)
		switch msg.NewState {
		case vo.LobbyView:
			cli.currentController = &lobby.Controller{
				Controller: cliVO.Controller{
					Model: lobby.Model{
						ViewModel: cliVO.ViewModel{
							Name: "Lobby",
						},
						ActiveMatchId: 0,
						AccountId:     msg.AccountId,
					},
					View: &lobby.View{},
				},
			}
			cli.ViewStateName = vo.LobbyView

			services.GetOpenMatches()
		case vo.JoinGameView:
			cli.isOpponent = true
			fallthrough
		case vo.CreateGameView:
			cli.matchId = msg.MatchId
			msg.AccountId = cli.account.ID //TODO fix this later

			resp := services.GetPlayerMatch(msg.MatchId)
			err := json.Unmarshal(resp.([]byte), &match)
			if err != nil {
				utils.Logger.Sugar().Error(err)
			}

			cli.GameMatch = match
			cli.createMatchView(match)
		}
	}

	cmds = append(cmds, cli.currentController.Update(msg, match))
	return cli, tea.Batch(cmds...)
}

func (cli *CLI) createMatchView(match *vo.GameMatch) {
	gameContainerModel := &container.Model{
		Tabs: []cliVO.Tab{
			{
				Title:    "Board",
				TabState: vo.BoardView,
			},
			{
				Title:    "Play",
				TabState: vo.PlayView,
			},
			{
				Title:    "Hand",
				TabState: vo.HandView,
			},
			{
				Title:    "Kitty",
				TabState: vo.KittyView,
			},
		},
		Match:       match,
		LocalPlayer: match.Players[0],
		ActiveTab:   0,
	}

	containerView := &container.View{
		ActiveTab: gameContainerModel.ActiveTab,
		Tabs:      gameContainerModel.Tabs,
	}

	cli.currentController = &gameContainer.Controller{
		Controller: &container.Controller{
			Controller: &cliVO.Controller{
				Model: gameContainerModel,
				View:  containerView,
			},
		},
	}

	ctrl := cli.currentController.(*gameContainer.Controller)
	ctrl.Init()
}

func GetMatchForId(msg vo.StateChangeMsg) (*vo.GameMatch, error) {
	var match *vo.GameMatch
	resp := services.GetPlayerMatch(msg.MatchId)
	err := json.Unmarshal(resp.([]byte), &match)
	if err != nil {
		utils.Logger.Sugar().Error(err)
	}
	return match, err
}
