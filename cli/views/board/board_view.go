package board

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/bardic/gocrib/cli/services"
	"github.com/bardic/gocrib/cli/styles"
	"github.com/bardic/gocrib/cli/utils"

	"github.com/bardic/gocrib/queries"
	"github.com/charmbracelet/bubbles/textinput"
)

type BoardView struct {
	cutInput             textinput.Model
	isLoading            bool //This should just be a state
	state                queries.Gamestate
	currentTurnsPlayerid int32
	localPlayer          *queries.Player
	players              []*queries.Player
	matchId              int32
}

var boardRowLen int = 50
var boardEndRowLen int = 31

func (v *BoardView) Init() {
	matchMsg := services.GetPlayerMatch(strconv.Itoa(int(v.matchId)))
	var match *queries.Match
	if err := json.Unmarshal(matchMsg.([]byte), &match); err != nil {
		return
	}

	v.cutInput = textinput.New()
	v.cutInput.Placeholder = "0"
	v.cutInput.CharLimit = 5
	v.cutInput.Width = 5
	v.isLoading = false
}

func (v *BoardView) Render() string {
	if v.isLoading {
		return "Loading..."
	}

	doc := strings.Builder{}
	viewBuilder := strings.Builder{}

	if v.state == queries.GamestateCutState && v.currentTurnsPlayerid != v.localPlayer.ID {
		v.cutInput.Focus()
		viewBuilder.WriteString(v.cutInput.View() + " \n")
	} else {
		viewBuilder.WriteString("\n")
	}

	//Row 1
	viewBuilder.WriteString(utils.DrawRow(v.players, boardRowLen, 0))
	//Row 2
	viewBuilder.WriteString(utils.DrawRow(v.players, boardRowLen, boardRowLen))
	//Row 3
	viewBuilder.WriteString(utils.DrawRow(v.players, boardEndRowLen, boardRowLen*2))

	doc.WriteString(styles.WindowStyle.Render(viewBuilder.String()))
	doc.WriteString(utils.BuildFooter())

	return doc.String()
}
