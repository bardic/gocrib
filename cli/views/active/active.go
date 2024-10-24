package active

import (
	"github.com/bardic/gocrib/cli/styles"
	"github.com/charmbracelet/lipgloss"
)

func ActiveView() string {
	return styles.WindowStyle.Align(lipgloss.Center, lipgloss.Center).Render("Active View")
}
