package main

import (
	"encoding/json"

	"cli/services"
	"cli/utils"
	"cli/view/container"
	"cli/view/gameContainer"
	"cli/view/lobby"
	cliVO "cli/vo"
	"queries"
	"vo"

	tea "github.com/charmbracelet/bubbletea"
)

func (cli *CLI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	cmds = append(cmds, cli.currentController.Update(msg))
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

			var matchDetails vo.MatchDetailsResponse
			msg := services.JoinMatch(cli.account.ID, msg.MatchId)
			err := json.Unmarshal(msg.([]byte), &matchDetails)

			if err != nil {

			}

		case vo.CreateGameView:
			cli.matchId = msg.MatchId
			msg.AccountId = cli.account.ID

			var match *vo.GameMatch
			resp := services.GetPlayerMatch(msg.MatchId)
			err := json.Unmarshal(resp.([]byte), &match)
			if err != nil {
				utils.Logger.Sugar().Error(err)
			}

			cli.createMatchView(match)
		}
	}

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
