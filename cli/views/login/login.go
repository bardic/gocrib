package login

import (
	"encoding/json"
	"strings"

	"github.com/bardic/gocrib/cli/services"
	"github.com/bardic/gocrib/cli/styles"
	"github.com/bardic/gocrib/cli/utils"
	"github.com/bardic/gocrib/queries"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type LoginView struct {
	LoginIdField      textinput.Model
	IsLoginIdFieldSet bool
}

func (view *LoginView) Init() {

	if view.IsLoginIdFieldSet {
		return
	}

	view.LoginIdField = textinput.New()
	view.LoginIdField.CharLimit = 20
	view.LoginIdField.Width = 30
	view.LoginIdField.Placeholder = "id"
	view.IsLoginIdFieldSet = true
}

func (v *LoginView) View() string {
	doc := strings.Builder{}
	doc.WriteString("Login \n")
	doc.WriteString(v.LoginIdField.View())

	return styles.WindowStyle.Align(lipgloss.Center, lipgloss.Center).Render(doc.String())
}

func (v *LoginView) ParseInput(msg tea.KeyMsg) tea.Msg {
	switch msg.String() {
	case "enter", "view_update":
		utils.Logger.Info("Enter")
		idStr := v.LoginIdField.Value()

		var accountDetails queries.Account
		msg := services.Login(idStr)
		json.Unmarshal(msg.([]byte), &accountDetails)

		return accountDetails
	}

	return nil
}

func (v *LoginView) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	v.LoginIdField.Focus()
	v.LoginIdField, cmd = v.LoginIdField.Update(msg)

	return cmd
}
