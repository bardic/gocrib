package active

import (
	"cli/styles"

	"github.com/charmbracelet/lipgloss"
)

func ActiveView() string {
	return styles.WindowStyle.Align(lipgloss.Center, lipgloss.Center).Render("Active View")
}
