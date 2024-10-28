package lobby

import (
	"encoding/json"
	"fmt"
	"strings"

	"cli/services"
	"cli/styles"
	"cli/utils"
	cliVO "cli/vo"
	"vo"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type View struct {
	AccountId        int32
	ActiveLandingTab int
	LobbyViewState   vo.ViewState
	LobbyTabNames    []string
	LobbyTable       table.Model
	IsLobbyTableSet  bool
	lobbyViewInitd   bool
	ActiveMatchId    int32
}

func (view *View) Init() {
	if view.lobbyViewInitd {
		return
	}

	view.lobbyViewInitd = true

	view.ActiveLandingTab = 0
	view.LobbyViewState = vo.OpenMatches
	view.LobbyTabNames = []string{"Open Matches", "Available Matches"}
}

func (view *View) Render() string {
	doc := strings.Builder{}

	renderedTabs := utils.RenderTabs([]cliVO.Tab{
		{
			Title:    "Lobby",
			TabState: vo.OpenMatches,
		},
		{
			Title:    "Active",
			TabState: vo.AvailableMatches,
		},
	}, view.ActiveLandingTab)

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	row = lipgloss.JoinHorizontal(lipgloss.Bottom, "─────────────────────────────────────────────────────────┐")
	doc.WriteString(row)
	doc.WriteString("\n")

	if !view.IsLobbyTableSet {
		t, err := getActiveView()
		if err != nil {
			return err.Error()
		}

		view.LobbyTable = t
		view.IsLobbyTableSet = true
	}

	doc.WriteString(styles.WindowStyle.Width(75).Height(12).Render(view.LobbyTable.View()))
	return doc.String()
}

func (view *View) Update(msg tea.Msg) tea.Cmd {
	view.Init()
	view.LobbyTable.Focus()

	updatedField, cmd := view.LobbyTable.Update(msg)
	view.LobbyTable = updatedField

	return cmd
}

func getActiveView() (table.Model, error) {
	columns := []table.Column{
		{Title: "Id", Width: 5},
		{Title: "Players", Width: 10},
		{Title: "Private", Width: 10},
		{Title: "Creation", Width: 20},
		{Title: "Turn", Width: 4},
		{Title: "State", Width: 5},
	}

	m := getOpenMatches()

	var matches []vo.GameMatch
	err := json.Unmarshal(m.([]byte), &matches)

	if err != nil {
		return table.Model{}, err
	}

	rows := []table.Row{}
	for _, m := range matches {

		rows = append(rows, table.Row{
			fmt.Sprintf("%v", m.ID),
			fmt.Sprintf("%v", m.Playerids),
			fmt.Sprintf("%v", m.Privatematch),
			m.Creationdate.Time.String(),
			fmt.Sprintf("%v", m.Currentplayerturn),
			fmt.Sprintf("%v", m.Gamestate),
		})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return t, nil
}

func getOpenMatches() tea.Msg {
	return services.GetOpenMatches()
}
