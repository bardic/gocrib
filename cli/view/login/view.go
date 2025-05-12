package login

import (
	"strings"

	"github.com/bardic/gocrib/cli/styles"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

type View struct {
	loginIDField      textinput.Model
	isLoginIDFieldSet bool
}

func NewLoginView() *View {
	v := &View{}
	v.Init()
	return v
}

func (view *View) Init() {
	if view.isLoginIDFieldSet {
		return
	}

	view.loginIDField = textinput.New()
	view.loginIDField.CharLimit = 20
	view.loginIDField.Width = 30
	view.loginIDField.Placeholder = "id"
	view.isLoginIDFieldSet = true
}

func (view *View) Render() string {
	doc := strings.Builder{}
	doc.WriteString("Login \n")
	doc.WriteString(view.loginIDField.View())

	return styles.WindowStyle.Align(lipgloss.Center, lipgloss.Center).Render(doc.String())
}

func (view *View) BuildHeader() string {
	return ""
}

func (view *View) BuildFooter() string {
	return ""
}
