package main

import (
	"github.com/bardic/gocrib/cli/styles"
)

func (m *CLI) View() string {
	return styles.ViewStyle.Render(m.currentController.Render(m.GameMatch))
}
