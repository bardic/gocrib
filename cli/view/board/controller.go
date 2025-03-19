package board

import (
	"github.com/bardic/gocrib/queries/queries"

	"github.com/bardic/gocrib/cli/services"
	cliVO "github.com/bardic/gocrib/cli/vo"
	"github.com/bardic/gocrib/vo"

	tea "github.com/charmbracelet/bubbletea"
)

type Controller struct {
	view  *View
	model *Model
}

func NewBoard(gameMatch *vo.GameMatch, player *vo.GamePlayer) *Controller {
	ctrl := &Controller{
		model: &Model{
			AccountId:   player.Accountid,
			GameMatchId: gameMatch.Match.ID,
			Gamestate:   gameMatch.Match.Gamestate,
		},
		view: &View{
			State:         queries.GamestateCut,
			Match:         gameMatch,
			LocalPlayerId: player.Accountid,
		},
	}

	ctrl.view.Init()

	return ctrl
}

func (ctrl *Controller) ShowInput() {
	ctrl.view.ShowCutInput()
}

func (ctrl *Controller) GetState() cliVO.ControllerState {
	return cliVO.BoardControllerState
}

func (ctrl *Controller) GetName() string {
	return "Board"
}

func (ctrl *Controller) Render(gameMatch *vo.GameMatch, gameDeck *vo.GameDeck) string {
	ctrl.model.Gamestate = gameMatch.Match.Gamestate
	return ctrl.view.Render()
}

func (ctrl *Controller) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd = []tea.Cmd{}
	ctrl.view.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (ctrl *Controller) ParseInput(msg tea.KeyMsg) tea.Msg {
	switch msg.String() {
	case "enter":
		return ctrl.Enter()
	}

	return msg
}

func (ctrl *Controller) Enter() tea.Msg {
	switch ctrl.model.Gamestate {
	case queries.GamestateCut:
		resp := services.CutDeck(*ctrl.model.GameMatchId, ctrl.view.CutInput.Value())
		return resp
	}

	return nil
}
