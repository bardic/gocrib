package main

import (
	"encoding/json"
	"strconv"

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
				NewState: vo.LobbyView,
			}
		})
	case vo.StateChangeMsg:
		switch msg.NewState {
		case vo.LobbyView:
			cli.currentController = &lobby.Controller{}
			cli.currentController.Init()
			cli.ViewStateName = vo.LobbyView

			services.GetOpenMatches()
		case vo.JoinGameView:

			fallthrough
		case vo.CreateGameView:
			cli.matchId = msg.MatchId

			var match *vo.GameMatch
			idstr := strconv.Itoa(msg.MatchId)
			resp := services.GetPlayerMatch(idstr)
			err := json.Unmarshal(resp.([]byte), &match)
			if err != nil {
				utils.Logger.Sugar().Error(err)
			}

			for _, player := range match.Players {
				if !player.Isready {
					services.PlayerReady(player.ID)
				}
			}

			resp = services.GetPlayerMatch(idstr)
			err = json.Unmarshal(resp.([]byte), &match)
			if err != nil {
				utils.Logger.Sugar().Error(err)
			}

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
	}

	return cli, tea.Batch(cmds...)
}
