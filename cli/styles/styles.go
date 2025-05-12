package styles

import "github.com/charmbracelet/lipgloss"

var (
	NormalBorderColour  = lipgloss.Color("10")
	FocusedBorderColour = lipgloss.Color("69")

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

	TabColour = lipgloss.
			AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}

	HelpColour = lipgloss.
			NewStyle().
			Foreground(lipgloss.Color("241"))
)

var (
	LineStyle = lipgloss.NewStyle().
			Width(10).
			Height(5)

	ModelStyle = LineStyle.
			BorderStyle(lipgloss.HiddenBorder())

	CardStyle = LineStyle.
			Align(lipgloss.Left, lipgloss.Top).
			BorderStyle(lipgloss.ThickBorder())

	SelectedStyle = CardStyle.
			BorderForeground(NormalBorderColour)

	SelectedFocusedStyle = CardStyle.
				BorderForeground(FocusedBorderColour)

	FocusedModelStyle = LineStyle.
				Align(lipgloss.Left, lipgloss.Top).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(FocusedBorderColour)

	ViewStyle = lipgloss.
			NewStyle().
			Padding(1, 2, 1, 2)

	InactiveTabStyle = lipgloss.NewStyle().
				Border(inactiveTabBorder, true).
				BorderForeground(TabColour).
				Padding(0, 1)

	ActiveTabStyle = InactiveTabStyle.
			Border(activeTabBorder, true)

	WindowStyle = lipgloss.NewStyle().
			BorderForeground(TabColour).
			Padding(1, 0).
			Border(lipgloss.NormalBorder()).
			BorderTop(true).
			Width(75).
			Height(12)

	WithStyleWithTabs = lipgloss.NewStyle().
				BorderForeground(TabColour).
				Padding(1, 0).
				Border(lipgloss.NormalBorder()).
				BorderTop(false).
				Width(75).
				Height(12)

	ActiveCard = WindowStyle.Align(lipgloss.Left, lipgloss.Top).Render("Active View")

	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
)

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}
