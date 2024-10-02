package views

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bardic/cribbagev2/cli/services"
	"github.com/bardic/cribbagev2/cli/styles"
	"github.com/bardic/cribbagev2/model"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var LobbyTable table.Model
var isLobbyTableSet bool

func getMatches() tea.Msg {
	return services.GetMatchesForPlayerId()
}

func getOpenMatches() tea.Msg {
	return services.GetOpenMatches()
}

func LobbyView(v ViewModel) (string, error) {

	doc := strings.Builder{}

	renderedTabs := renderTabs(v.LandingTabs, v.ActiveLandingTab)

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	row = lipgloss.JoinHorizontal(lipgloss.Bottom, "───────────────────────────────────────────────────────────────────┐")
	doc.WriteString(row)
	doc.WriteString("\n")

	switch v.ViewState {
	case model.LoginView:

	case model.LobbyView:
		if isLobbyTableSet {
			break
		}

		t, err := getActiveView()
		if err != nil {
			return "", err
		}

		LobbyTable = t
		isLobbyTableSet = true
	}

	doc.WriteString(styles.WindowStyle.Width(100).Render(LobbyTable.View()))
	return doc.String(), nil
}

func getActiveView() (table.Model, error) {
	columns := []table.Column{
		{Title: "Id", Width: 5},
		{Title: "Players", Width: 10},
		{Title: "Private", Width: 10},
		{Title: "Creation", Width: 20},
		{Title: "Turn", Width: 4},
		{Title: "Last Play", Width: 20},
		{Title: "State", Width: 5},
	}

	m := getOpenMatches()

	var matchesStr string
	err := json.Unmarshal(m.([]byte), &matchesStr)

	if err != nil {
		return table.Model{}, err
	}

	var matches []model.Match
	err = json.Unmarshal([]byte(matchesStr), &matches)

	if err != nil  {
		return table.Model{}, err
	}

	rows := []table.Row{}
	for _, m := range matches {
		l := len(m.TurnPassTimestamps)
		lastTurnTimestamp := ""
		if l > 0 {
			lastTurnTimestamp = m.TurnPassTimestamps[l-1]
		}

		rows = append(rows, table.Row{
			fmt.Sprintf("%v", m.Id),
			fmt.Sprintf("%v", m.PlayerIds),
			fmt.Sprintf("%v", m.PrivateMatch),
			m.CreationDate.String(),
			fmt.Sprintf("%v", m.CurrentPlayerTurn),
			lastTurnTimestamp,
			fmt.Sprintf("%v", m.GameState),
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
