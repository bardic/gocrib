package styles

import "github.com/charmbracelet/lipgloss"

var (
	ModelStyle = lipgloss.NewStyle().
			Width(10).
			Height(5).
			Align(lipgloss.Center, lipgloss.Center).
			BorderStyle(lipgloss.HiddenBorder())

	SelectedStyle = lipgloss.NewStyle().
			Width(10).
			Height(5).
			Align(lipgloss.Center, lipgloss.Center).
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(lipgloss.Color("10"))

	SelectedFocusedStyle = lipgloss.NewStyle().
				Width(10).
				Height(5).
				Align(lipgloss.Center, lipgloss.Center).
				BorderStyle(lipgloss.ThickBorder()).
				BorderForeground(lipgloss.Color("69"))

	FocusedModelStyle = lipgloss.NewStyle().
				Width(10).
				Height(5).
				Align(lipgloss.Center, lipgloss.Center).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("69"))

	HelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")

	ViewStyle        = lipgloss.NewStyle().Padding(1, 2, 1, 2)
	highlightColor   = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	InactiveTabStyle = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 1)
	ActiveTabStyle   = InactiveTabStyle.Border(activeTabBorder, true)
	WindowStyle      = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(2, 0).Border(lipgloss.NormalBorder())
)

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}
