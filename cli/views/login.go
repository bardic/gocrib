package views

import (
	"fmt"
	"strings"

	"github.com/bardic/cribbagev2/cli/styles"
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
	fmt.Println("Enter")
}
