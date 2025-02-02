package board

import (
	"encoding/json"
	"strings"

	"github.com/bardic/gocrib/vo"

	"github.com/bardic/gocrib/cli/services"
	"github.com/bardic/gocrib/cli/utils"

	"github.com/bardic/gocrib/queries/queries"

	"github.com/charmbracelet/bubbles/textinput"
)

type View struct {
	cutInput      textinput.Model
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

	view.cutInput = textinput.New()
	view.cutInput.Placeholder = "0"
	view.cutInput.CharLimit = 5
	view.cutInput.Focus()
	view.cutInput.Width = 5
	view.isLoading = false
}

func (view *View) Render(hand []int) string {
	if view.isLoading {
		return "Loading..."
	}

	doc := strings.Builder{}
	viewBuilder := strings.Builder{}

	if view.State == queries.GamestateCut {
		viewBuilder.WriteString(view.cutInput.View() + " \n")
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

func (view *View) BuildHeader() string {
	return ""
}

func (view *View) BuildFooter() string {
	f := utils.BuildCommonFooter(
		view.Match.Currentplayerturn,
		view.LocalPlayerId,
		view.Match.ID,
		view.Match.Gamestate,
	)
	return f
}
