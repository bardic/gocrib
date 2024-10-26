package main

import (
	"encoding/json"
	"strconv"

	"cli/services"
	"cli/utils"
	"cli/view/container"
	"cli/view/lobby"
	cliVO "cli/vo"
	"queries"
	"vo"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	cmds = append(cmds, m.currentController.Update(msg))
	switch msg := msg.(type) {
	case queries.Account:
		m.account = &msg
		cmds = append(cmds, func() tea.Msg {
			return vo.StateChangeMsg{
				NewState: vo.LobbyView,
			}
		})
	case vo.StateChangeMsg:
		switch msg.NewState {
		case vo.LobbyView:
			m.currentController = &lobby.Controller{}
			m.currentController.Init()
			m.ViewStateName = vo.LobbyView

			services.GetOpenMatches()
		case vo.JoinGameView:

			fallthrough
		case vo.CreateGameView:
			m.matchId = msg.MatchId

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

			containerModel := &container.Model{
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
				Match: match,
			}

			containerView := &container.View{
				ActiveTab: containerModel.ActiveTab,
				Tabs:      containerModel.Tabs,
			}

			m.currentController = &container.Controller{
				Controller: cliVO.Controller{
					Model: containerModel,
					View:  containerView,
				},
			}

			ctrl := m.currentController.(*container.Controller)
			ctrl.Init()

		}
	}

	return m, tea.Batch(cmds...)
}
