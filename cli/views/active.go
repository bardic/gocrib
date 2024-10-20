package views

import (
	"github.com/bardic/gocrib/cli/styles"
	"github.com/charmbracelet/lipgloss"
)

func ActiveView() string {
	return styles.ScreenStyle.Width(100).Align(lipgloss.Center, lipgloss.Center).Render("Active View")
}
