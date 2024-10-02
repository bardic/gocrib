package views

import (
	"strings"

	"github.com/bardic/cribbagev2/cli/styles"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

var LoginIdField textinput.Model
var isLoginIdFieldSet bool

func LoginView() string {
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
