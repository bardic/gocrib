package container

import (
	"encoding/json"
	"time"

	"github.com/bardic/gocrib/cli/services"
	"github.com/bardic/gocrib/cli/styles"
	"github.com/bardic/gocrib/cli/utils"
	"github.com/bardic/gocrib/cli/view/board"
	"github.com/bardic/gocrib/cli/view/card"
	"github.com/bardic/gocrib/cli/view/game"
	cliVO "github.com/bardic/gocrib/cli/vo"
	"github.com/bardic/gocrib/vo"

	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
)

type Controller struct {
	*game.Controller
	timer        timer.Model
	timerStarted bool
}

func (ctrl *Controller) Init() {

}

func NewController(match *vo.GameMatch, player *vo.GamePlayer) *Controller {
	model := NewModel(match, player)
	view := NewView(model)
	ctrl := &Controller{
		Controller: &game.Controller{
			Model: model,
			View:  view,
		},
	}

	ctrl.ChangeTab(0)
	return ctrl
}

func (ctrl *Controller) GetState() cliVO.ControllerState {
	return cliVO.LoginControllerState
}

func (ctrl *Controller) Render(gamematch *vo.GameMatch) string {
	containerModel := ctrl.Model.(*Model)

	cardIds := []int{}

	for _, card := range containerModel.GetPlayer().Hand {
		cardIds = append(cardIds, *card.Cardid)
	}

	containerHeader := ctrl.GetView().Render(cardIds)
	viewRender := containerModel.GetSubcontroller().Render(containerModel.ViewModel.Gamematch)

	return containerHeader + "\n" + styles.WindowStyle.Render(viewRender)
}

func (ctrl *Controller) Update(msg tea.Msg, gameMatch *vo.GameMatch) tea.Cmd {
	var cmds []tea.Cmd
	containerModel := ctrl.GetModel().(*Model)
	subView := containerModel.GetSubcontroller()

	if gameMatch != nil {
		cmd := subView.Update(msg, gameMatch)
		cmds = append(cmds, cmd)
	}

	if !ctrl.timerStarted {
		ctrl.timer = timer.NewWithInterval(time.Hour, time.Second*1)
		cmd := ctrl.timer.Init()
		cmds = append(cmds, cmd)
		ctrl.timerStarted = true
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		resp := ctrl.ParseInput(msg)

		if resp == nil {
			break
		}

		switch r := resp.(type) {
		case vo.ChangeTabMsg:
			cmds = append(cmds, func() tea.Msg {
				return r
			})
		}
	case timer.TickMsg: // Polling update
		var cmd tea.Cmd
		var gameMatch *vo.GameMatch
		ctrl.timer, cmd = ctrl.timer.Update(msg)

		resp := services.GetMatchById(ctrl.GetModel().GetMatch().ID)
		err := json.Unmarshal(resp.([]byte), &gameMatch)

		if err != nil {
			utils.Logger.Sugar().Error(err)
		}

		ctrl.Controller.GetModel().(*Model).ViewModel.Gamematch = gameMatch

		cmds = append(cmds, cmd)
	case vo.ChangeTabMsg:
		ctrl.ChangeTab(msg.TabIndex)
	}

	return tea.Batch(cmds...)
}

func (ctrl *Controller) ParseInput(msg tea.KeyMsg) tea.Msg {
	containerModel := ctrl.GetModel().(*Model)
	containerView := ctrl.GetView().(*View)

	switch msg.String() {
	case "ctrl+c", "q":
		return tea.Quit()
	case "tab":
		containerView.ActiveTab = containerView.ActiveTab + 1
		if containerView.ActiveTab >= len(containerView.Tabs) {
			containerView.ActiveTab = 0
		}
		containerModel.State = containerView.Tabs[containerView.ActiveTab].TabState
		return vo.ChangeTabMsg{
			TabIndex: containerView.ActiveTab,
		}

	case "shift+tab":
		containerView.ActiveTab = containerView.ActiveTab - 1

		if containerView.ActiveTab < 0 {
			containerView.ActiveTab = len(containerView.Tabs) - 1
		}

		containerModel.State = containerView.Tabs[containerView.ActiveTab].TabState
		return vo.ChangeTabMsg{
			TabIndex: containerView.ActiveTab,
		}
	default:
		containerModel.Subcontroller.ParseInput(msg)
	}

	return msg
}

func (ctrl *Controller) ChangeTab(tabIndex int) {
	containerModel := ctrl.Controller.GetModel().(*Model)

	switch tabIndex {
	case 0:
		containerModel.Subcontroller = &board.Controller{
			Controller: &game.Controller{
				Model: &board.Model{
					ViewModel: cliVO.ViewModel{
						Name: "Game",
					},
					GameMatch: containerModel.ViewModel.Gamematch,
				},
				View: &board.View{
					Match:         containerModel.ViewModel.Gamematch,
					LocalPlayerId: containerModel.ViewModel.GetPlayer().ID,
					State:         containerModel.ViewModel.Gamematch.Gamestate,
				},
			},
		}
	case 1:
		containerModel.Subcontroller = card.NewController(
			"Play",
			containerModel.GetMatch(),
			containerModel.GetPlayer(),
		)
	case 2:

		containerModel.Subcontroller = card.NewController(
			"Hand",
			containerModel.GetMatch(),
			containerModel.GetPlayer(),
		)
	case 3:

		containerModel.Subcontroller = card.NewController(
			"Kitty",
			containerModel.GetMatch(),
			containerModel.GetPlayer(),
		)
	}
}
