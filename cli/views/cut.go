package views

import (
	"github.com/bardic/cribbagev2/cli/styles"
	"github.com/charmbracelet/lipgloss"
)

func CutView() string {
	return styles.WindowStyle.Width(100).Align(lipgloss.Center, lipgloss.Center).Render("Cut View")
}
