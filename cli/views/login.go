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
}

var LoginIdField textinput.Model
var isLoginIdFieldSet bool

func (s LoginView) View() string {
	doc := strings.Builder{}
	if !isLoginIdFieldSet {
		LoginIdField = textinput.New()
		LoginIdField.CharLimit = 20
		LoginIdField.Width = 30
		LoginIdField.Placeholder = "id"
		isLoginIdFieldSet = true
	}

	doc.WriteString("Login \n")
	doc.WriteString(LoginIdField.View())

	return styles.ScreenStyle.Width(100).Align(lipgloss.Center, lipgloss.Center).Render(doc.String())
}

func (s LoginView) Enter() tea.Msg {
	utils.Logger.Info("Enter")
	idStr := LoginIdField.Value()
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return tea.Quit
	}

	state.AccountId = id
	state.ViewStateName = model.LobbyView
	return services.Login()
}
