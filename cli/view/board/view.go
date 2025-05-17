package board

import (
	"encoding/json"
	"strings"

	"github.com/bardic/gocrib/cli/services"
	"github.com/bardic/gocrib/cli/utils"
	"github.com/bardic/gocrib/vo"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type View struct {
	CutInput      textinput.Model
	isLoading     bool
	State         string
	LocalPlayerID int
	Match         *vo.Match
}

var (
	boardRowLen    = 50
	boardEndRowLen = 31
)

func (view *View) Init() {
	matchMsg := services.GetMatchByID(view.Match.ID)
	var match *vo.Match
	if err := json.Unmarshal(matchMsg.([]byte), &match); err != nil {
		return
	}

	view.ShowCutInput()
}

func (view *View) ShowCutInput() {
	view.CutInput = textinput.New()
	view.CutInput.Placeholder = "0"
	view.CutInput.Cursor.Blink = true
	view.CutInput.CharLimit = 5
	view.CutInput.Focus()
	view.CutInput.Width = 5
	view.isLoading = false
}

func (view *View) Render() string {
	if view.isLoading {
		return "Loading..."
	}

	doc := strings.Builder{}
	viewBuilder := strings.Builder{}

	if view.State == "Cut" {
		viewBuilder.WriteString(view.CutInput.View() + " \n")
	} else {
		viewBuilder.WriteString("\n")
	}

	// Row 1
	viewBuilder.WriteString(utils.DrawRow(view.Match.Players, boardRowLen, 0))
	// Row 2
	viewBuilder.WriteString(utils.DrawRow(view.Match.Players, boardRowLen, boardRowLen))
	// Row 3
	viewBuilder.WriteString(utils.DrawRow(view.Match.Players, boardEndRowLen, boardRowLen*2))

	doc.WriteString(viewBuilder.String())
	doc.WriteString(view.BuildFooter())

	return doc.String()
}

func (view *View) Update(msg tea.Msg) {
	if view.isLoading {
		return
	}

	view.CutInput.Focus()
	view.CutInput, _ = view.CutInput.Update(msg)
}

func (view *View) BuildHeader() string {
	return ""
}

func (view *View) BuildFooter() string {
	f := utils.BuildCommonFooter(
		view.Match,
	)
	return f
}
