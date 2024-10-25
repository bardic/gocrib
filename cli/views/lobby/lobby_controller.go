package lobby

import (
	"encoding/json"
	"strconv"

	"cli/services"
	"cli/utils"
	"cli/views"
	"model"
	"queries"

	tea "github.com/charmbracelet/bubbletea"
)

type LobbyController struct {
	views.Controller
}

func (gc *LobbyController) GetState() views.ControllerState {
	return views.LobbyControllerState
}

func (gc *LobbyController) Init() {
	gc.Model = LobbyModel{}
	gc.View = &LobbyView{}
}
func (gc *LobbyController) Render() string {
	return gc.View.Render()
}

func (v *LobbyController) ParseInput(msg tea.KeyMsg) tea.Msg {
	lobbyView := v.View.(*LobbyView)
	lobbyModel := v.Model.(LobbyModel)

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

		var matchDetails model.MatchDetailsResponse
		msg := services.JoinMatch(int(player.ID), id)
		err = json.Unmarshal(msg.([]byte), &matchDetails)

		if err != nil {
			return tea.Quit
		}

		return model.StateChangeMsg{
			NewState: model.JoinGameView,
			MatchId:  matchDetails.MatchId,
		}
	case "n":
		match := utils.CreateGame(lobbyModel.AccountId)

		return model.StateChangeMsg{
			NewState: model.CreateGameView,
			MatchId:  match.MatchId,
		}
	case "tab":

		lobbyView.ActiveLandingTab = lobbyView.ActiveLandingTab + 1

		switch lobbyView.ActiveLandingTab {
		case 0:
			lobbyView.LobbyViewState = model.OpenMatches
		case 1:
			lobbyView.LobbyViewState = model.AvailableMatches
		}
	case "shift+tab":
		lobbyView.ActiveLandingTab = lobbyView.ActiveLandingTab - 1

		switch lobbyView.ActiveLandingTab {
		case 0:
			lobbyView.LobbyViewState = model.OpenMatches
		case 1:
			lobbyView.LobbyViewState = model.AvailableMatches
		}
	}

	return nil
}

func (v *LobbyController) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	lobbyView := v.View.(*LobbyView)
	lobbyView.LobbyTable.Focus()

	updatedField, cmd := lobbyView.LobbyTable.Update(msg)
	lobbyView.LobbyTable = updatedField

	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg: //User input
		resp := v.ParseInput(msg)

		if resp == nil {
			break
		}

		cmds = append(cmds, func() tea.Msg {
			return resp
		})

	}

	return tea.Batch(cmds...)
}
