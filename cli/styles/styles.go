package styles

import "github.com/charmbracelet/lipgloss"

var (
	ValueTopLeft    = lipgloss.NewStyle().Align(lipgloss.Left, lipgloss.Top)
	ValueBottoRight = lipgloss.NewStyle().Align(lipgloss.Right, lipgloss.Bottom)

	ModelStyle = lipgloss.NewStyle().
			Width(10).
			Height(5).
		// Align(lipgloss.Left, lipgloss.Top).
		BorderStyle(lipgloss.HiddenBorder())

	SelectedStyle = lipgloss.NewStyle().
			Width(10).
			Height(5).
			Align(lipgloss.Left, lipgloss.Top).
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(lipgloss.Color("10"))

	SelectedFocusedStyle = lipgloss.NewStyle().
				Width(10).
				Height(5).
				Align(lipgloss.Left, lipgloss.Top).
				BorderStyle(lipgloss.ThickBorder()).
				BorderForeground(lipgloss.Color("69"))

	FocusedModelStyle = lipgloss.NewStyle().
				Width(10).
				Height(5).
				Align(lipgloss.Left, lipgloss.Top).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("69"))

	HelpStyle = lipgloss.
			NewStyle().
			Foreground(lipgloss.Color("241"))

	ViewStyle = lipgloss.
			NewStyle().
			Padding(1, 2, 1, 2)

	highlightColor = lipgloss.
			AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}

	InactiveTabStyle = lipgloss.NewStyle().
				Border(inactiveTabBorder, true).
				BorderForeground(highlightColor).
				Padding(0, 1)

	ActiveTabStyle = InactiveTabStyle.
			Border(activeTabBorder, true)

	WindowStyle = lipgloss.NewStyle().
			BorderForeground(highlightColor).
			Padding(1, 0).
			Border(lipgloss.NormalBorder()).
			BorderTop(false).
			Width(75).
			Height(12)

	ActiveCard = WindowStyle.Align(lipgloss.Left, lipgloss.Top).Render("Active View")

	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")

	Player1 = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
	Player2 = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
	Player3 = lipgloss.NewStyle().Foreground(lipgloss.Color("#0000FF"))
	Player4 = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00"))

	PlayerStyles = []lipgloss.Style{
		Player1,
		Player2,
		Player3,
		Player4,
	}
)

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}
