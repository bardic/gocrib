package lobby

import (
	"encoding/json"
	"strconv"

	"cli/services"
	"cli/utils"
	cliVO "cli/vo"
	"queries"
	"vo"

	tea "github.com/charmbracelet/bubbletea"
)

type Controller struct {
	cliVO.Controller
}

func (ctrl *Controller) GetState() cliVO.ControllerState {
	return cliVO.LobbyControllerState
}

func (ctrl *Controller) Init() {
	ctrl.Model = Model{}
	ctrl.View = &View{}
}
func (ctrl *Controller) Render() string {
	return ctrl.View.Render()
}

func (ctrl *Controller) ParseInput(msg tea.KeyMsg) tea.Msg {
	lobbyView := ctrl.View.(*View)
	lobbyModel := ctrl.Model.(Model)

	switch msg.String() {
	case "enter", "view_update":
		utils.Logger.Info("Enter")
		idStr := lobbyView.LobbyTable.SelectedRow()[0]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return tea.Quit
		}

		accountMsg := services.PostPlayer(int32(id))

		var player queries.Player
		err = json.Unmarshal(accountMsg.([]byte), &player)

		if err != nil {
			return tea.Quit
		}

		var matchDetails vo.MatchDetailsResponse
		msg := services.JoinMatch(int(player.ID), id)
		err = json.Unmarshal(msg.([]byte), &matchDetails)

		if err != nil {
			return tea.Quit
		}

		return vo.StateChangeMsg{
			NewState: vo.JoinGameView,
			MatchId:  matchDetails.MatchId,
		}
	case "n":
		match := utils.CreateGame(lobbyModel.AccountId)

		return vo.StateChangeMsg{
			NewState: vo.CreateGameView,
			MatchId:  match.MatchId,
		}
	case "tab":

		lobbyView.ActiveLandingTab = lobbyView.ActiveLandingTab + 1

		switch lobbyView.ActiveLandingTab {
		case 0:
			lobbyView.LobbyViewState = vo.OpenMatches
		case 1:
			lobbyView.LobbyViewState = vo.AvailableMatches
		}
	case "shift+tab":
		lobbyView.ActiveLandingTab = lobbyView.ActiveLandingTab - 1

		switch lobbyView.ActiveLandingTab {
		case 0:
			lobbyView.LobbyViewState = vo.OpenMatches
		case 1:
			lobbyView.LobbyViewState = vo.AvailableMatches
		}
	}

	return nil
}

func (ctrl *Controller) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	lobbyView := ctrl.View.(*View)
	lobbyView.LobbyTable.Focus()

	updatedField, cmd := lobbyView.LobbyTable.Update(msg)
	lobbyView.LobbyTable = updatedField

	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg: //User input
		resp := ctrl.ParseInput(msg)

		if resp == nil {
			break
		}

		cmds = append(cmds, func() tea.Msg {
			return resp
		})

	}

	return tea.Batch(cmds...)
}
