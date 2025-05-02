package login

import (
	"encoding/json"

	"github.com/bardic/gocrib/queries/queries"
	"github.com/bardic/gocrib/vo"

	"github.com/bardic/gocrib/cli/services"
	logger "github.com/bardic/gocrib/cli/utils/log"
	cliVO "github.com/bardic/gocrib/cli/vo"

	tea "github.com/charmbracelet/bubbletea"
)

type Controller struct {
	model *Model
	view  *View
}

func NewLogin() *Controller {
	ctrl := &Controller{
		model: &Model{},
		view:  NewLoginView(),
	}

	return ctrl
}

func (ctrl *Controller) GetState() cliVO.ControllerState {
	return cliVO.LoginControllerState
}

func (ctrl *Controller) GetName() string {
	return "Login"
}

func (ctrl *Controller) Render() string {
	return ctrl.view.Render()
}

func (ctrl *Controller) ParseInput(msg tea.KeyMsg) tea.Msg {
	l := logger.Get()
	defer l.Sync()

	switch msg.String() {
	case "ctrl+c", "q":
		return tea.Quit()
	case "enter", "view_update":
		l.Sugar().Info("Enter")
		idStr := ctrl.view.loginIdField.Value()

		var accountDetails queries.Account
		msg := services.Login(idStr)
		json.Unmarshal(msg.([]byte), &accountDetails)

		return vo.StateChangeMsg{
			NewState:  vo.LobbyView,
			AccountId: accountDetails.ID,
		}

		// return accountDetails
	}

	return nil
}

func (ctrl *Controller) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	ctrl.view.loginIdField.Focus()
	ctrl.view.loginIdField, cmd = ctrl.view.loginIdField.Update(msg)

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
