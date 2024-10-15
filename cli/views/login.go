package views

import (
	"strconv"
	"strings"

	"github.com/bardic/gocrib/cli/services"
	"github.com/bardic/gocrib/cli/state"
	"github.com/bardic/gocrib/cli/styles"
	"github.com/bardic/gocrib/cli/utils"
	"github.com/bardic/gocrib/model"
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

	return styles.ScreenStyle.Width(100).Align(lipgloss.Center, lipgloss.Center).Render(doc.String())
}

func (v *LoginView) Enter() tea.Msg {
	utils.Logger.Info("Enter")
	idStr := v.LoginIdField.Value()
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return tea.Quit
	}

	state.AccountId = id
	state.ViewStateName = model.LobbyView
	return services.Login()
}

func (v *LoginView) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	v.LoginIdField, cmd = v.LoginIdField.Update(msg)

	return cmd
}
