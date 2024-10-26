package main

import (
	"cli/styles"
)

func (m *CLI) View() string {
	return styles.ViewStyle.Render(m.currentController.Render())
}
