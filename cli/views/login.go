package views

import (
	"strings"

	"github.com/bardic/gocrib/cli/styles"
	"github.com/bardic/gocrib/cli/utils"
	"github.com/charmbracelet/bubbles/textinput"
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

func (s LoginView) Enter() {
	utils.Logger.Info("Enter")
}
