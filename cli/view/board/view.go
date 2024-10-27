package board

import (
	"encoding/json"
	"strconv"
	"strings"
	"vo"

	"cli/services"
	"cli/utils"

	"queries"

	"github.com/charmbracelet/bubbles/textinput"
)

type View struct {
	cutInput      textinput.Model
	isLoading     bool //This should just be a state
	state         queries.Gamestate
	localPlayerId int32
	match         *vo.GameMatch
}

var boardRowLen int = 50
var boardEndRowLen int = 31

func (view *View) Init() {
	matchMsg := services.GetPlayerMatch(strconv.Itoa(int(view.match.ID)))
	var match *queries.Match
	if err := json.Unmarshal(matchMsg.([]byte), &match); err != nil {
		return
	}

	view.cutInput = textinput.New()
	view.cutInput.Placeholder = "0"
	view.cutInput.CharLimit = 5
	view.cutInput.Width = 5
	view.isLoading = false
}

func (view *View) Render() string {
	if view.isLoading {
		return "Loading..."
	}

	doc := strings.Builder{}
	viewBuilder := strings.Builder{}

	if view.state == queries.GamestateCutState && view.match.Currentplayerturn != view.localPlayerId {
		view.cutInput.Focus()
		viewBuilder.WriteString(view.cutInput.View() + " \n")
	} else {
		viewBuilder.WriteString("\n")
	}

	//Row 1
	viewBuilder.WriteString(utils.DrawRow(view.match.Players, boardRowLen, 0))
	//Row 2
	viewBuilder.WriteString(utils.DrawRow(view.match.Players, boardRowLen, boardRowLen))
	//Row 3
	viewBuilder.WriteString(utils.DrawRow(view.match.Players, boardEndRowLen, boardRowLen*2))

	doc.WriteString(viewBuilder.String())
	doc.WriteString(utils.BuildFooter())

	return doc.String()
}
