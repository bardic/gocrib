package views

import (
	"github.com/bardic/cribbagev2/cli/styles"
	"github.com/charmbracelet/lipgloss"
)

func LobbyView() string {
	return styles.ScreenStyle.Width(100).Align(lipgloss.Center, lipgloss.Center).Render("Lobby View")
}
