package board

import (
	"github.com/bardic/gocrib/cli/services"
	"github.com/bardic/gocrib/cli/views"
	"github.com/bardic/gocrib/queries"
	tea "github.com/charmbracelet/bubbletea"
)

type BoardController struct {
	views.Controller
}

func (ctrl *BoardController) GetState() views.ControllerState {
	return views.BoardControllerState
}

func (ctrl *BoardController) Init() {
	boardModel := ctrl.Model.(BoardModel)
	ctrl.View = &BoardView{
		matchId:              boardModel.GameMatch.ID,
		players:              boardModel.GameMatch.Players,
		state:                queries.GamestateCutState,
		localPlayer:          boardModel.GameMatch.Players[0],
		currentTurnsPlayerid: 0,
	}

	ctrl.View.Init()
}

func (ctrl *BoardController) Render() string {
	return ctrl.View.Render()
}

func (ctrl *BoardController) Update(msg tea.Msg) tea.Cmd {
	gameView := ctrl.View.(*BoardView)
	var cmd tea.Cmd
	gameView.cutInput.Focus()
	gameView.cutInput, cmd = gameView.cutInput.Update(msg)
	return cmd
}

func (ctrl *BoardController) ParseInput(msg tea.KeyMsg) tea.Msg {
	switch msg.String() {
	case "ctrl+c", "q":
		return tea.Quit()
	case "enter", "view_update":
		return ctrl.Enter()

	}
	return nil
}

func (ctrl *BoardController) Enter() tea.Msg {
	boardView := ctrl.View.(*BoardView)
	boardModel := ctrl.Model.(*BoardModel)
	switch boardModel.GameMatch.Gamestate {
	case queries.GamestateCutState:
		boardModel.CutIndex = boardView.cutInput.Value()
		resp := services.CutDeck(boardModel.Account.ID, boardModel.GameMatch.ID, boardModel.CutIndex)
		return resp
	}

	return nil
}
