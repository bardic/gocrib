package board

import (
	"encoding/json"
	"strings"

	"github.com/bardic/gocrib/vo"

	"github.com/bardic/gocrib/cli/services"
	"github.com/bardic/gocrib/cli/utils"

	"github.com/bardic/gocrib/queries/queries"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type View struct {
	CutInput      textinput.Model
	isLoading     bool
	State         queries.Gamestate
	LocalPlayerId *int
	Match         *vo.GameMatch
}

var boardRowLen int = 50
var boardEndRowLen int = 31

func (view *View) Init() {
	matchMsg := services.GetMatchById(view.Match.ID)
	var match *queries.Match
	if err := json.Unmarshal(matchMsg.([]byte), &match); err != nil {
		return
	}
}

func (view *View) ShowCutInput() {
	view.CutInput = textinput.New()
	view.CutInput.Placeholder = "0"
	view.CutInput.CharLimit = 5
	view.CutInput.Focus()
	view.CutInput.Width = 5
	view.isLoading = false

}

func (view *View) Render(hand []int) string {
	if view.isLoading {
		return "Loading..."
	}

	doc := strings.Builder{}
	viewBuilder := strings.Builder{}

	if view.State == queries.GamestateCut {
		viewBuilder.WriteString(view.CutInput.View() + " \n")
	} else {
		viewBuilder.WriteString("\n")
	}

	//Row 1
	viewBuilder.WriteString(utils.DrawRow(view.Match.Players, boardRowLen, 0))
	//Row 2
	viewBuilder.WriteString(utils.DrawRow(view.Match.Players, boardRowLen, boardRowLen))
	//Row 3
	viewBuilder.WriteString(utils.DrawRow(view.Match.Players, boardEndRowLen, boardRowLen*2))

	doc.WriteString(viewBuilder.String())
	doc.WriteString(view.BuildFooter())

	return doc.String()
}

func (view *View) Update(msg tea.Msg) {
	view.CutInput, _ = view.CutInput.Update(msg)

}

func (view *View) BuildHeader() string {
	return ""
}

func (view *View) BuildFooter() string {
	p := utils.GetPlayerForAccountId(view.LocalPlayerId, view.Match)
	f := utils.BuildCommonFooter(
		view.Match,
		p,
	)
	return f
}
