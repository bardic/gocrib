package views

import (
	tea "github.com/charmbracelet/bubbletea"
)

type IViewState interface {
	ParseInput(tea.KeyMsg) tea.Msg
	Enter() tea.Msg
	View() string
	Update(msg tea.Msg) tea.Cmd
}
