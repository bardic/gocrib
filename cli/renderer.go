package main

import (
	"cli/styles"
)

func (m *AppModel) View() string {
	return styles.ViewStyle.Render(m.currentController.Render())
}
