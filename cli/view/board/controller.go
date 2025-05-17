package board

import (
	"github.com/bardic/gocrib/cli/services"
	cliVO "github.com/bardic/gocrib/cli/vo"
	"github.com/bardic/gocrib/vo"

	tea "github.com/charmbracelet/bubbletea"
)

type Controller struct {
	view  *View
	model *Model
}

func NewBoard(gameMatch *vo.Match, player *vo.Player) *Controller {
	ctrl := &Controller{
		model: &Model{
			AccountID:   player.Accountid,
			GameMatchID: gameMatch.ID,
			Gamestate:   gameMatch.Gamestate,
		},
		view: &View{
			State:         "Cut",
			Match:         gameMatch,
			LocalPlayerID: player.Accountid,
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

func (ctrl *Controller) Render(gameMatch *vo.Match, gameDeck *vo.Deck) string {
	ctrl.model.Gamestate = gameMatch.Gamestate
	return ctrl.view.Render()
}

func (ctrl *Controller) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	cmds := []tea.Cmd{}
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
	case "Cut":
		resp := services.CutDeck(ctrl.model.GameMatchID, ctrl.view.CutInput.Value())
		return resp
	}

	return nil
}
