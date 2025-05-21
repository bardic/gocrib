package lobby

import (
	"encoding/json"
	"strconv"

	"github.com/bardic/gocrib/cli/services"
	cliVO "github.com/bardic/gocrib/cli/vo"
	logger "github.com/bardic/gocrib/log"
	"github.com/bardic/gocrib/vo"
	tea "github.com/charmbracelet/bubbletea"
)

type Controller struct {
	model *Model
	view  *View
}

func NewLobby(accountID int) *Controller {
	return &Controller{
		model: &Model{
			playerAccountID: accountID,
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

		var match *vo.Match
		msg := services.JoinMatch(lobbyModel.playerAccountID, id)
		err = json.Unmarshal(msg.([]byte), &match)
		if err != nil {
			return tea.Quit
		}

		return cliVO.ChangeState{
			NewState: "Join",
			MatchID:  id,
		}
	case "n":
		match := CreateGame(lobbyModel.playerAccountID)

		return cliVO.ChangeState{
			NewState:  "Create",
			AccountID: lobbyModel.playerAccountID,
			MatchID:   match.ID,
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

func CreateGame(accountID int) *vo.Match {
	newMatch := services.PostPlayerMatch(accountID).([]byte)

	var match *vo.Match
	err := json.Unmarshal(newMatch, &match)
	if err != nil {
		return nil
	}

	return match
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
