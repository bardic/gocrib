package login

import (
	"strings"

	"github.com/bardic/gocrib/cli/styles"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

type View struct {
	loginIdField      textinput.Model
	isLoginIdFieldSet bool
}

func (view *View) Init() {
	if view.isLoginIdFieldSet {
		return
	}

	view.loginIdField = textinput.New()
	view.loginIdField.CharLimit = 20
	view.loginIdField.Width = 30
	view.loginIdField.Placeholder = "id"
	view.isLoginIdFieldSet = true
}

func (view *View) Render(hand []int32) string {
	doc := strings.Builder{}
	doc.WriteString("Login \n")
	doc.WriteString(view.loginIdField.View())

	return styles.WindowStyle.Align(lipgloss.Center, lipgloss.Center).Render(doc.String())
}

func (view *View) BuildHeader() string {
	return ""
}

func (view *View) BuildFooter() string {
	return ""
}
