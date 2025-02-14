package container

import (
	"encoding/json"
	"time"

	"github.com/bardic/gocrib/cli/services"
	"github.com/bardic/gocrib/cli/styles"
	"github.com/bardic/gocrib/cli/utils"
	"github.com/bardic/gocrib/cli/view/board"
	"github.com/bardic/gocrib/cli/view/card"
	cliVO "github.com/bardic/gocrib/cli/vo"
	"github.com/bardic/gocrib/vo"

	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
)

type Controller struct {
	View         *View
	Model        *Model
	timer        timer.Model
	timerStarted bool
	tabIndex     int
	tabs         map[int]cliVO.IController
}

func (ctrl *Controller) Init() {

}

func NewController(match *vo.GameMatch, player *vo.GamePlayer) *Controller {
	tabs := createTabs(match, player)
	ctrl := &Controller{
		Model: NewModel(match, player),
		View:  NewView(0, tabs),
	}

	ctrl.tabs = tabs

	ctrl.ChangeTab(vo.ChangeTabMsg{
		TabIndex: 0,
	})
	return ctrl
}

func (ctrl *Controller) GetModel() cliVO.IModel {
	return ctrl.Model
}

func (ctrl *Controller) GetView() cliVO.IView {
	return ctrl.View
}

func (ctrl *Controller) GetName() string {
	return "Container"
}

func (ctrl *Controller) GetState() cliVO.ControllerState {
	return cliVO.LoginControllerState
}

func (ctrl *Controller) Render(gamematch *vo.GameMatch) string {
	viewRender := ctrl.Model.Subcontroller.Render(gamematch)

	return ctrl.View.Render() + "\n" + styles.WindowStyle.Render(viewRender)
}

func (ctrl *Controller) Update(msg tea.Msg, gameMatch *vo.GameMatch) tea.Cmd {
	var cmds []tea.Cmd

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

		resp := services.GetMatchById(ctrl.Model.GetMatch().ID)
		err := json.Unmarshal(resp.([]byte), &gameMatch)

		if err != nil {
			utils.Logger.Sugar().Error(err)
		}

		ctrl.Model.SetMatch(gameMatch)
		// ctrl.Model.Subcontroller.Render(gameMatch)
		cmds = append(cmds, cmd)
	case vo.ChangeTabMsg:
		ctrl.ChangeTab(msg)
		ctrl.tabs[ctrl.tabIndex].Update(msg, ctrl.Model.GetMatch())
	}

	return tea.Batch(cmds...)
}

func (ctrl *Controller) ParseInput(msg tea.KeyMsg) tea.Msg {
	switch msg.String() {
	case "ctrl+c", "q":
		return tea.Quit()
	case "tab":
		ctrl.View.ActiveTab = ctrl.View.ActiveTab + 1
		if ctrl.View.ActiveTab >= len(ctrl.View.Tabs) {
			ctrl.View.ActiveTab = 0
		}
		return vo.ChangeTabMsg{
			TabIndex: ctrl.View.ActiveTab,
		}

	case "shift+tab":
		ctrl.View.ActiveTab = ctrl.View.ActiveTab - 1

		if ctrl.View.ActiveTab < 0 {
			ctrl.View.ActiveTab = len(ctrl.View.Tabs) - 1
		}

		return vo.ChangeTabMsg{
			TabIndex: ctrl.View.ActiveTab,
		}
	default:
		ctrl.Model.Subcontroller.ParseInput(msg)
	}

	return msg
}

func (ctrl *Controller) ChangeTab(msg tea.Msg) {
	tabIndex := msg.(vo.ChangeTabMsg).TabIndex
	if ctrl.tabs == nil {
		ctrl.tabs = map[int]cliVO.IController{}
	}

	ctrl.tabIndex = tabIndex

	val, ok := ctrl.tabs[tabIndex]
	if ok {
		ctrl.Model.Subcontroller = val
		val.Update(msg, ctrl.Model.GetMatch())
		return
	}

}

func createTabs(gameMatch *vo.GameMatch, player *vo.GamePlayer) map[int]cliVO.IController {
	return map[int]cliVO.IController{
		0: board.NewBoard(
			gameMatch,
			player,
		),
		1: card.NewController(
			"Play",
			gameMatch,
			player,
		),
		2: card.NewController(
			"Hand",
			gameMatch,
			player,
		),
		3: card.NewController(
			"Kitty",
			gameMatch,
			player,
		),
	}
}
