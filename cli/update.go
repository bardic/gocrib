package main

import (
	"encoding/json"
	"strconv"

	"cli/services"
	"cli/utils"
	"cli/views"
	"cli/views/container"
	"cli/views/lobby"
	"model"
	"queries"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	cmds = append(cmds, m.currentController.Update(msg))
	switch msg := msg.(type) {
	case queries.Account:
		m.account = &msg
		cmds = append(cmds, func() tea.Msg {
			return model.StateChangeMsg{
				NewState: model.LobbyView,
			}
		})
	case model.StateChangeMsg:
		switch msg.NewState {
		case model.LobbyView:
			m.currentController = &lobby.LobbyController{}
			m.currentController.Init()
			m.ViewStateName = model.LobbyView

			services.GetOpenMatches()
		case model.JoinGameView:

			fallthrough
		case model.CreateGameView:
			m.matchId = msg.MatchId

			var match *model.GameMatch
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

			containerModel := &container.ContainerModel{
				Tabs: []views.Tab{
					{
						Title:    "Board",
						TabState: model.BoardView,
					},
					{
						Title:    "Play",
						TabState: model.PlayView,
					},
					{
						Title:    "Hand",
						TabState: model.HandView,
					},
					{
						Title:    "Kitty",
						TabState: model.KittyView,
					},
				},
				Match: match,
			}

			containerView := &container.ContainerView{
				ActiveTab: containerModel.ActiveTab,
				Tabs:      containerModel.Tabs,
			}

			m.currentController = &container.ContainerController{
				Controller: views.Controller{
					Model: containerModel,
					View:  containerView,
				},
			}

			ctrl := m.currentController.(*container.ContainerController)
			ctrl.Init()

		}
	}

	return m, tea.Batch(cmds...)
}
