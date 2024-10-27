package board

import (
	"cli/services"
	cliVO "cli/vo"
	"queries"

	tea "github.com/charmbracelet/bubbletea"
)

type Controller struct {
	cliVO.Controller
}

func (ctrl *Controller) GetState() cliVO.ControllerState {
	return cliVO.BoardControllerState
}

func (ctrl *Controller) Init() {
	boardModel := ctrl.Model.(Model)
	ctrl.View = &View{
		state:         queries.GamestateCutState,
		match:         boardModel.GameMatch,
		localPlayerId: boardModel.LocalPlayerId,
	}

	ctrl.View.Init()
}

func (ctrl *Controller) Render() string {
	return ctrl.View.Render()
}

func (ctrl *Controller) Update(msg tea.Msg) tea.Cmd {
	gameView := ctrl.View.(*View)
	var cmd tea.Cmd
	gameView.cutInput.Focus()
	gameView.cutInput, cmd = gameView.cutInput.Update(msg)
	return cmd
}

func (ctrl *Controller) ParseInput(msg tea.KeyMsg) tea.Msg {
	switch msg.String() {
	case "ctrl+c", "q":
		return tea.Quit()
	case "enter", "view_update":
		return ctrl.Enter()

	}
	return nil
}

func (ctrl *Controller) Enter() tea.Msg {
	boardView := ctrl.View.(*View)
	boardModel := ctrl.Model.(*Model)
	switch boardModel.GameMatch.Gamestate {
	case queries.GamestateCutState:
		boardModel.CutIndex = boardView.cutInput.Value()
		resp := services.CutDeck(boardModel.Account.ID, boardModel.GameMatch.ID, boardModel.CutIndex)
		return resp
	}

	return nil
}
