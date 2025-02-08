package login

import (
	"encoding/json"

	"github.com/bardic/gocrib/queries/queries"

	"github.com/bardic/gocrib/cli/services"
	"github.com/bardic/gocrib/cli/utils"
	"github.com/bardic/gocrib/cli/view/container"
	cliVO "github.com/bardic/gocrib/cli/vo"
	"github.com/bardic/gocrib/vo"

	tea "github.com/charmbracelet/bubbletea"
)

type Controller struct {
	cliVO.GameController
}

func (ctrl *Controller) GetModel() cliVO.IModel {
	return &container.Model{}
}

func New() *Controller {
	ctrl := &Controller{}

	ctrl.Model = LoginModel{
		cliVO.ViewModel{
			Name: "Login",
		},
	}

	ctrl.View = &View{}
	ctrl.View.Init()

	return ctrl
}

func (ctrl *Controller) GetState() cliVO.ControllerState {
	return cliVO.LoginControllerState
}

func (ctrl *Controller) Init() {
	ctrl.Model = LoginModel{
		cliVO.ViewModel{
			Name: "Login",
		},
	}

	ctrl.View = &View{}
	ctrl.View.Init()
}

func (ctrl *Controller) Render(gamematch *vo.GameMatch) string {
	return ctrl.View.Render(nil)
}

func (ctrl *Controller) ParseInput(msg tea.KeyMsg) tea.Msg {
	loginView := ctrl.View.(*View)

	switch msg.String() {
	case "ctrl+c", "q":
		return tea.Quit()
	case "enter", "view_update":
		utils.Logger.Info("Enter")
		idStr := loginView.loginIdField.Value()

		var accountDetails queries.Account
		msg := services.Login(idStr)
		json.Unmarshal(msg.([]byte), &accountDetails)

		return accountDetails
	}

	return nil
}

func (ctrl *Controller) Update(msg tea.Msg, gameMatch *vo.GameMatch) tea.Cmd {
	loginView := ctrl.View.(*View)

	var cmd tea.Cmd
	var cmds []tea.Cmd

	loginView.loginIdField.Focus()
	loginView.loginIdField, cmd = loginView.loginIdField.Update(msg)

	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg: //User input
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
