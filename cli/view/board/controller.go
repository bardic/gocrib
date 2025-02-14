package board

import (
	"github.com/bardic/gocrib/queries/queries"

	"github.com/bardic/gocrib/cli/services"
	"github.com/bardic/gocrib/cli/view/game"
	cliVO "github.com/bardic/gocrib/cli/vo"
	"github.com/bardic/gocrib/vo"

	tea "github.com/charmbracelet/bubbletea"
)

type Controller struct {
	*game.Controller
}

func NewBoard(gameMatch *vo.GameMatch, player *vo.GamePlayer) *Controller {
	ctrl := &Controller{
		Controller: &game.Controller{
			Model: &Model{
				ViewModel: cliVO.ViewModel{
					Gamematch: gameMatch,
					Name:      "Board",
					AccountId: player.Accountid,
				},
				AccountId: player.Accountid,
			},
			View: &View{
				State:         queries.GamestateCut,
				Match:         gameMatch,
				LocalPlayerId: player.Accountid,
			},
		},
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

func (ctrl *Controller) GetName() string {
	return "Board"
}

func (ctrl *Controller) Render(gameMatch *vo.GameMatch) string {
	ids := []int{}

	for _, card := range gameMatch.Players[0].Hand {
		ids = append(ids, *card.Cardid)
	}

	return ctrl.View.Render()
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
	}

	return msg
}

func (ctrl *Controller) Enter() tea.Msg {
	boardView := ctrl.View.(*View)
	boardModel := ctrl.Model.(*Model)
	switch boardModel.ViewModel.Gamematch.Gamestate {
	case queries.GamestateCut:
		boardModel.CutIndex = boardView.CutInput.Value()
		resp := services.CutDeck(*boardModel.ViewModel.Gamematch.ID, boardModel.CutIndex)
		return resp
	case queries.GamestatePlay:
		resp := services.PutPlay(boardModel.AccountId, boardModel.ViewModel.Gamematch.ID, vo.HandModifier{})
		return resp
	}

	return nil
}
