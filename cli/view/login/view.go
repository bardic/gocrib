package login

import (
	"strings"

	"cli/styles"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

type LoginView struct {
	loginIdField      textinput.Model
	isLoginIdFieldSet bool
}

func (view *LoginView) Init() {
	if view.isLoginIdFieldSet {
		return
	}

	view.loginIdField = textinput.New()
	view.loginIdField.CharLimit = 20
	view.loginIdField.Width = 30
	view.loginIdField.Placeholder = "id"
	view.isLoginIdFieldSet = true
}

func (view *LoginView) Render() string {
	doc := strings.Builder{}
	doc.WriteString("Login \n")
	doc.WriteString(view.loginIdField.View())

	return styles.WindowStyle.Align(lipgloss.Center, lipgloss.Center).Render(doc.String())
}
