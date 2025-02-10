package board

import (
	"github.com/bardic/gocrib/queries/queries"

	"github.com/bardic/gocrib/cli/services"
	"github.com/bardic/gocrib/cli/utils"
	"github.com/bardic/gocrib/cli/view/game"
	cliVO "github.com/bardic/gocrib/cli/vo"
	"github.com/bardic/gocrib/vo"

	tea "github.com/charmbracelet/bubbletea"
)

type Controller struct {
	*game.Controller
}

func NewBoard() *Controller {
	ctrl := &Controller{
		Controller: &game.Controller{},
	}
	boardModel := ctrl.Model.(*Model)

	ctrl.View = &View{
		State:         queries.GamestateCut,
		Match:         boardModel.GameMatch,
		LocalPlayerId: utils.GetPlayerForAccountId(boardModel.Account.ID, boardModel.GameMatch).ID,
	}

	ctrl.View.(*View).Init()

	return ctrl
}

func (ctrl *Controller) ShowInput() {
	ctrl.View.(*View).ShowCutInput()
}

func (ctrl *Controller) GetState() cliVO.ControllerState {
	return cliVO.BoardControllerState
}

func (ctrl *Controller) Render(gameMatch *vo.GameMatch) string {
	ids := []int{}

	for _, card := range gameMatch.Players[0].Hand {
		ids = append(ids, *card.Cardid)
	}

	return ctrl.View.Render(ids)
}

func (ctrl *Controller) Update(msg tea.Msg, gameMatch *vo.GameMatch) tea.Cmd {
	var cmd tea.Cmd
	ctrl.View.(*View).Update(msg)
	return cmd
}

func (ctrl *Controller) ParseInput(msg tea.KeyMsg) tea.Msg {
	switch msg.String() {
	case "enter":
		return ctrl.Enter()
		// default:
		// 	ctrl.View.(*View).Update(msg)
	}

	return msg
}

func (ctrl *Controller) Enter() tea.Msg {
	boardView := ctrl.View.(*View)
	boardModel := ctrl.Model.(*Model)
	switch boardModel.GameMatch.Gamestate {
	case queries.GamestateCut:
		boardModel.CutIndex = boardView.CutInput.Value()
		resp := services.CutDeck(*boardModel.GameMatch.ID, boardModel.CutIndex)
		playerId := utils.GetPlayerForAccountId(boardModel.Account.ID, boardModel.GameMatch).ID
		services.PlayerReady(playerId, boardModel.GameMatch.ID)
		return resp
	}

	return nil
}
