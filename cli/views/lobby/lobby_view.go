package lobby

import (
	"encoding/json"
	"fmt"
	"strings"

	"cli/services"
	"cli/styles"
	"cli/utils"
	"cli/views"
	"model"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jackc/pgx/v5/pgtype"
)

type LobbyView struct {
	AccountId        int32
	ActiveLandingTab int
	LobbyViewState   model.ViewState
	LobbyTabNames    []string
	LobbyTable       table.Model
	IsLobbyTableSet  bool
	lobbyViewInitd   bool
	ActiveMatchId    int32
}

func (v *LobbyView) Init() {
	if v.lobbyViewInitd {
		return
	}

	v.lobbyViewInitd = true

	v.ActiveLandingTab = 0
	v.LobbyViewState = model.OpenMatches
	v.LobbyTabNames = []string{"Open Matches", "Available Matches"}
}

func (v *LobbyView) Render() string {
	doc := strings.Builder{}

	renderedTabs := utils.RenderTabs([]views.Tab{
		{
			Title:    "Lobby",
			TabState: model.OpenMatches,
		},
		{
			Title:    "Active",
			TabState: model.AvailableMatches,
		},
	}, v.ActiveLandingTab)

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	row = lipgloss.JoinHorizontal(lipgloss.Bottom, "────────────────────────────────────────────────────────────────┐")
	doc.WriteString(row)
	doc.WriteString("\n")

	if !v.IsLobbyTableSet {
		t, err := getActiveView()
		if err != nil {
			return err.Error()
		}

		v.LobbyTable = t
		v.IsLobbyTableSet = true
	}

	doc.WriteString(styles.WindowStyle.Width(100).Render(v.LobbyTable.View()))
	return doc.String()
}

func (v *LobbyView) Update(msg tea.Msg) tea.Cmd {
	v.Init()
	v.LobbyTable.Focus()

	updatedField, cmd := v.LobbyTable.Update(msg)
	v.LobbyTable = updatedField

	return cmd
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

	var matches []model.GameMatch
	err := json.Unmarshal(m.([]byte), &matches)

	if err != nil {
		return table.Model{}, err
	}

	rows := []table.Row{}
	for _, m := range matches {
		l := len(m.Turnpasstimestamps)
		var lastTurnTimestamp pgtype.Timestamptz
		if l > 0 {
			lastTurnTimestamp = m.Turnpasstimestamps[l-1]
		}

		rows = append(rows, table.Row{
			fmt.Sprintf("%v", m.ID),
			fmt.Sprintf("%v", m.Playerids),
			fmt.Sprintf("%v", m.Privatematch),
			m.Creationdate.Time.String(),
			fmt.Sprintf("%v", m.Currentplayerturn),
			lastTurnTimestamp.Time.String(),
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
