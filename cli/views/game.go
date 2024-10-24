package views

import (
	"encoding/json"
	"slices"
	"strconv"
	"strings"

	"github.com/bardic/gocrib/cli/services"
	"github.com/bardic/gocrib/cli/styles"
	"github.com/bardic/gocrib/cli/utils"

	"github.com/bardic/gocrib/model"
	"github.com/bardic/gocrib/queries"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type GameView struct {
	Initd          bool
	MatchId        int32
	Account        *queries.Account
	CutIndex       string
	HighlightedIds []int32
	GameTabNames   []string
	GameViewState  model.GameViewState
	ActiveSlot     int
	ActiveTab      int
	HighlighedId   int
	CutInput       textinput.Model
	Deck           *model.GameDeck
	DeckId         int32
	Hand           []queries.Card
	Kitty          []queries.Card
	Play           []queries.Card
	GameMatch      *model.GameMatch
	LocalPlayer    *queries.Player
}

var boardRowLen int = 50
var boardEndRowLen int = 31

func (v *GameView) Init() {
	if v.Initd {
		return
	}

	v.Hand = []queries.Card{}
	v.Kitty = []queries.Card{}
	v.Play = []queries.Card{}

	matchMsg := services.GetPlayerMatch(strconv.Itoa(int(v.MatchId)))
	var match *queries.Match
	if err := json.Unmarshal(matchMsg.([]byte), &match); err != nil {
		return
	}

	v.Initd = true
	v.GameTabNames = []string{"Board", "Play", "Hand", "Kitty"}
	v.CutInput = textinput.New()
	v.CutInput.Placeholder = "0"
	v.CutInput.CharLimit = 5
	v.CutInput.Width = 5
}

func (v *GameView) View() string {
	if v.GameMatch == nil {
		return "Loading..."
	}

	doc := strings.Builder{}
	renderedTabs := renderTabs(v.GameTabNames, v.ActiveTab)
	doc.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...))
	doc.WriteString(lipgloss.JoinHorizontal(lipgloss.Bottom, "───────────────────────────────────────────────────────────────────┐"))
	viewBuilder := strings.Builder{}

	switch v.GameViewState {

	case model.BoardView:
		player, err := utils.GetPlayerById(v.Account.ID, v.GameMatch.Players)
		if err != nil {
			utils.Logger.Sugar().Error(err)
		}
		if v.GameMatch.Gamestate == queries.GamestateCutState && v.GameMatch.Currentplayerturn != player.ID {
			v.CutInput.Focus()
			viewBuilder.WriteString(v.CutInput.View() + " \n")
		} else {
			viewBuilder.WriteString("\n")
		}

		//Row 1
		viewBuilder.WriteString(utils.DrawRow(v.GameMatch.Players, boardRowLen, 0))
		//Row 2
		viewBuilder.WriteString(utils.DrawRow(v.GameMatch.Players, boardRowLen, boardRowLen))
		//Row 3
		viewBuilder.WriteString(utils.DrawRow(v.GameMatch.Players, boardEndRowLen, boardRowLen*2))

	case model.PlayView:
		hand := v.createHandView()
		playerView := PlayerView{
			HandModel: hand,
		}
		viewBuilder.WriteString(playerView.View())
	case model.HandView:
		hand := v.createHandView()
		handView := HandView{
			HandModel: hand,
		}

		if v.GameMatch.Gamestate == queries.GamestateWaitingState {
			viewBuilder.WriteString("Waiting to be dealt")
		} else {
			viewBuilder.WriteString(handView.View())
		}
	case model.KittyView:
		hand := v.createHandView()
		kittyView := KittyView{
			HandModel: hand,
		}

		if len(hand.player.Kitty) == 0 {
			viewBuilder.WriteString("Empty Kitty")
		} else {
			viewBuilder.WriteString(kittyView.View())
		}
	}

	doc.WriteString(styles.WindowStyle.Render(viewBuilder.String()))
	doc.WriteString(utils.BuildFooter())
	return doc.String()
}

func (v *GameView) createHandView() HandModel {

	p := utils.GetPlayerForAccountId(v.Account.ID, v.GameMatch)

	hand := HandModel{
		currentTurnPlayerId: v.GameMatch.Currentplayerturn,
		selectedCardId:      v.HighlighedId,
		selectedCardIds:     v.HighlightedIds,
		cards:               p.Play,
		deck:                v.Deck,
		player:              p,
		account:             v.Account,
	}

	return hand
}

func (v *GameView) Enter() tea.Msg {
	switch v.GameMatch.Gamestate {
	case queries.GamestateCutState:
		v.CutIndex = v.CutInput.Value()
		resp := services.CutDeck(v.Account.ID, v.MatchId, v.CutIndex)
		return resp
	case queries.GamestateDiscardState:
		p, err := utils.GetPlayerById(v.Account.ID, v.GameMatch.Players)

		if err != nil {
			utils.Logger.Sugar().Error(err)
		}

		services.PutKitty(model.HandModifier{
			MatchId:  v.GameMatch.ID,
			PlayerId: p.ID,
			CardIds:  v.HighlightedIds,
		})

		v.HighlightedIds = []int32{}
	case queries.GamestatePlayState:
		p, err := utils.GetPlayerById(v.Account.ID, v.GameMatch.Players)

		if err != nil {
			utils.Logger.Sugar().Error(err)
		}

		services.PutPlay(model.HandModifier{
			MatchId:  v.GameMatch.ID,
			PlayerId: p.ID,
			CardIds:  v.HighlightedIds,
		})

		v.HighlightedIds = []int32{}
	}

	return nil
}

func (v *GameView) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	v.CutInput.Focus()
	v.CutInput, cmd = v.CutInput.Update(msg)
	return cmd
}

func (v *GameView) ParseInput(msg tea.KeyMsg) tea.Msg {
	switch msg.String() {
	case "ctrl+c", "q":
		return tea.Quit()
	case "enter", "view_update":
		return v.Enter()
	case " ":
		cards := utils.GetVisibleCards(v.ActiveTab, v.GameMatch.Players[0])

		if len(cards) == 0 {
			return nil
		}

		idx := slices.Index(v.HighlightedIds, cards[v.HighlighedId])
		if idx > -1 {
			v.HighlightedIds = slices.Delete(v.HighlightedIds, 0, 1)
		} else {
			v.HighlightedIds = append(v.HighlightedIds, cards[v.HighlighedId])
		}
	case "tab":
		v.ActiveTab = v.ActiveTab + 1

		switch v.ActiveTab {
		case 0:
			v.GameViewState = model.BoardView
		case 1:
			v.GameViewState = model.PlayView
		case 2:
			v.GameViewState = model.HandView
		case 3:
			v.GameViewState = model.KittyView
		}

	case "shift+tab":
		v.ActiveTab = v.ActiveTab - 1

		switch v.ActiveTab {
		case 0:
			v.GameViewState = model.BoardView
		case 1:
			v.GameViewState = model.PlayView
		case 2:
			v.GameViewState = model.HandView
		case 3:
			v.GameViewState = model.KittyView
		}

	case "right":
		v.ActiveSlot++

		cards := utils.GetVisibleCards(v.ActiveTab, v.GameMatch.Players[0])

		if v.ActiveSlot > len(cards)-1 {
			v.ActiveSlot = 0
		}

		v.HighlighedId = v.ActiveSlot
	case "left":

		v.ActiveSlot--

		cards := utils.GetVisibleCards(v.ActiveTab, v.GameMatch.Players[0])

		if v.ActiveSlot < 0 {
			v.ActiveSlot = len(cards) - 1
		}

		v.HighlighedId = v.ActiveSlot
	}

	return nil
}
