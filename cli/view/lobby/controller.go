package lobby

import (
	"encoding/json"
	"strconv"

	"github.com/bardic/gocrib/cli/services"
	logger "github.com/bardic/gocrib/cli/utils/log"
	cliVO "github.com/bardic/gocrib/cli/vo"
	"github.com/bardic/gocrib/vo"

	tea "github.com/charmbracelet/bubbletea"
)

type Controller struct {
	model *Model
	view  *View
}

func NewLobby(msg vo.StateChangeMsg) *Controller {
	return &Controller{
		model: &Model{
			playerAccountID: msg.AccountID,
		},
		view: &View{
			ActiveLandingTab: 0,
		},
	}
}

func (ctrl *Controller) GetState() cliVO.ControllerState {
	return cliVO.LobbyControllerState
}

func (ctrl *Controller) GetName() string {
	return "Login"
}

func (ctrl *Controller) Init() {
}

func (ctrl *Controller) Render() string {
	return ctrl.view.Render()
}

func (ctrl *Controller) ParseInput(msg tea.KeyMsg) tea.Msg {
	l := logger.Get()
	defer l.Sync()
	lobbyView := ctrl.view
	lobbyModel := ctrl.model

	switch msg.String() {
	case "ctrl+c", "q":
		return tea.Quit()
	case "enter", "view_update":
		l.Sugar().Info("Enter")
		if len(lobbyView.LobbyTable.Rows()) == 0 {
			return nil
		}
		idStr := lobbyView.LobbyTable.SelectedRow()[0]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return tea.Quit
		}

		var matchDetails vo.MatchDetailsResponse
		msg := services.JoinMatch(*lobbyModel.playerAccountID, id)
		err = json.Unmarshal(msg.([]byte), &matchDetails)
		if err != nil {
			return tea.Quit
		}

		return vo.StateChangeMsg{
			NewState: vo.JoinGameView,
			MatchID:  &id,
		}
	case "n":
		match := CreateGame(lobbyModel.playerAccountID)

		return vo.StateChangeMsg{
			NewState:  vo.CreateGameView,
			AccountID: lobbyModel.playerAccountID,
			MatchID:   match.MatchID,
		}
	case "tab":

		lobbyView.ActiveLandingTab++

		switch lobbyView.ActiveLandingTab {
		case 0:
			lobbyView.LobbyViewState = vo.OpenMatches
		case 1:
			lobbyView.LobbyViewState = vo.AvailableMatches
		}
	case "shift+tab":
		lobbyView.ActiveLandingTab--

		switch lobbyView.ActiveLandingTab {
		case 0:
			lobbyView.LobbyViewState = vo.OpenMatches
		case 1:
			lobbyView.LobbyViewState = vo.AvailableMatches
		}
	}

	return nil
}

func CreateGame(accountID *int) vo.MatchDetailsResponse {
	newMatch := services.PostPlayerMatch(accountID).([]byte)

	var matchDetails vo.MatchDetailsResponse
	err := json.Unmarshal(newMatch, &matchDetails)
	if err != nil {
		return vo.MatchDetailsResponse{}
	}

	return matchDetails
}

func (ctrl *Controller) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	lobbyView := ctrl.view
	lobbyView.LobbyTable.Focus()

	updatedField, cmd := lobbyView.LobbyTable.Update(msg)
	lobbyView.LobbyTable = updatedField

	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg: // User input
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
