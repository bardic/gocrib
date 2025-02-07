package board

import (
	"github.com/bardic/gocrib/queries/queries"

	"github.com/bardic/gocrib/cli/services"
	cliVO "github.com/bardic/gocrib/cli/vo"
	"github.com/bardic/gocrib/vo"

	tea "github.com/charmbracelet/bubbletea"
)

type Controller struct {
	*cliVO.Controller
}

func (ctrl *Controller) GetState() cliVO.ControllerState {
	return cliVO.BoardControllerState
}

func (ctrl *Controller) Init() {
	boardModel := ctrl.Model.(*Model)
	ctrl.View = &View{
		State:         queries.GamestateCut,
		Match:         boardModel.GameMatch,
		LocalPlayerId: boardModel.LocalPlayerId,
	}

	ctrl.View.Init()

}

func (ctrl *Controller) Render(gameMatch *vo.GameMatch) string {
	ids := []int{}

	for _, card := range gameMatch.Players[0].Hand {
		ids = append(ids, *card.Cardid)
	}

	return ctrl.View.Render(ids)
}

func (ctrl *Controller) Update(msg tea.Msg, gameMatch *vo.GameMatch) tea.Cmd {
	gameView := ctrl.View.(*View)
	var cmd tea.Cmd

	gameView.cutInput, cmd = gameView.cutInput.Update(msg)
	gameView.cutInput.Focus()

	return cmd
}

func (ctrl *Controller) ParseInput(msg tea.KeyMsg) tea.Msg {
	switch msg.String() {
	case "enter":
		return ctrl.Enter()
	}

	return msg
}

func (ctrl *Controller) Enter() tea.Msg {
	boardView := ctrl.View.(*View)
	boardModel := ctrl.Model.(*Model)
	switch boardModel.GameMatch.Gamestate {
	case queries.GamestateCut:
		boardModel.CutIndex = boardView.cutInput.Value()
		resp := services.CutDeck(*boardModel.GameMatch.ID, boardModel.CutIndex)
		services.PlayerReady(boardModel.LocalPlayerId, boardModel.GameMatch.ID)
		return resp
	}

	return nil
}
