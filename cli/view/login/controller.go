package login

import (
	"encoding/json"

	"cli/services"
	"cli/utils"
	cliVO "cli/vo"
	"queries"

	tea "github.com/charmbracelet/bubbletea"
)

type LoginController struct {
	cliVO.Controller
}

func (ctrl *LoginController) GetState() cliVO.ControllerState {
	return cliVO.LoginControllerState
}

func (ctrl *LoginController) Init() {
	ctrl.Model = LoginModel{
		cliVO.ViewModel{
			Name: "Login",
		},
	}

	ctrl.View = &LoginView{}
	ctrl.View.Init()
}

func (ctrl *LoginController) Render() string {
	return ctrl.View.Render()
}

func (ctrl *LoginController) ParseInput(msg tea.KeyMsg) tea.Msg {
	loginView := ctrl.View.(*LoginView)

	switch msg.String() {
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

func (ctrl *LoginController) Update(msg tea.Msg) tea.Cmd {
	loginView := ctrl.View.(*LoginView)

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
