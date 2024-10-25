package main

import (
	"github.com/bardic/gocrib/cli/styles"
)

func (m *AppModel) View() string {
	return styles.ViewStyle.Render(m.currentController.Render())
}
