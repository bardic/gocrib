package views

import (
	"encoding/json"
	"strings"

	"github.com/bardic/gocrib/cli/services"
	"github.com/bardic/gocrib/cli/state"
	"github.com/bardic/gocrib/cli/styles"
	"github.com/bardic/gocrib/cli/utils"
	"github.com/bardic/gocrib/model"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type GameView struct {
	GameState      model.GameState
	PreviousState  model.GameState
	HighlightedIds []int
	Cards          []model.Card
	GameTabNames   []string
	GameViewState  model.GameViewState
	ActiveSlot     int
	ActiveTab      int
	HighlighedId   int
	CutInput       textinput.Model
	gameViewInitd  bool
	DeckId         int
	Hand           []model.Card
	Kitty          []model.Card
	Play           []model.Card
	GameMatch      *model.GameMatch
}

var focusedModelStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("69"))

func (v *GameView) Init() {
	if v.gameViewInitd {
		return
	}

	v.gameViewInitd = true
	v.GameTabNames = []string{"Board", "Play", "Hand", "Kitty"}
	v.CutInput = textinput.New()
	v.CutInput.Placeholder = "0"
	v.CutInput.CharLimit = 5
	v.CutInput.Width = 5

	deckByte := services.GetDeckById(v.DeckId).([]byte)
	var deck model.GameDeck
	json.Unmarshal(deckByte, &deck)
	state.ActiveDeck = &deck
}

func (v *GameView) View() string {
	doc := strings.Builder{}

	renderedTabs := renderTabs(v.GameTabNames, v.ActiveTab)

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	row = lipgloss.JoinHorizontal(lipgloss.Bottom, "───────────────────────────────────────────────────────────────────┐")
	doc.WriteString(row)
	doc.WriteString("\n")

	var view string
	switch v.GameViewState {
	case model.BoardView:
		if state.ActiveMatch != nil && state.ActiveMatch.GameState == model.CutState {
			v.CutInput.Focus()
			view = v.CutInput.View() + " \n"
		} else {
			view = "\n"
		}

		view = view + `
	○•○○○ ○○○○○ ○○○○○ ○○○○○ ○○○○○ ○○○○○ ○○○○○ ○○○○○ ○○○○○ ○○○○○		
	○•○○○ ○○○○○ ○○○○○ ○○○○○ ○○○○○ ○○○○○ ○○○○○ ○○○○○ ○○○○○ ○○○○○		
	----- ----- ----- ----- ----- ----- ----- ----- ----- -----		
	○○○○○ ○○○○○ ○○○○○ ○○○○○ ○○○○○ ○○○○○ ○○○○○ ○○○○○ ○○○○○ ○○○○○		
	○○○○○ ○○○○○ ○○○○○ ○○○○○ ○○○○○ ○○○○○ ○○○○○ ○○○○○ ○○○○○ ○○○○○		
	----- ----- ----- ----- ----- ----- ----- ----- ----- -----		
	○○○○○ ○○○○○ ○○○○○ ○○○○○ ○○○○○ ○○○○○ ○
	○○○○○ ○○○○○ ○○○○○ ○○○○○ ○○○○○ ○○○○○ ○
`

		s := lipgloss.JoinHorizontal(lipgloss.Top, focusedModelStyle.Render(view), focusedModelStyle.Render("59"))
		view = s
	case model.PlayView:
		view = HandView(v.HighlighedId, v.HighlightedIds, v.Cards)
	case model.HandView:
		if v.GameState == 0 {
			view = "Waiting to be dealt"
		} else {
			view = HandView(v.HighlighedId, v.HighlightedIds, v.Cards)
		}
	case model.KittyView:
		if len(v.Cards) == 0 {
			view = "Empty Kitty"
		} else {
			view = HandView(v.HighlighedId, v.HighlightedIds, v.Cards)
		}
	}

	doc.WriteString(styles.WindowStyle.Width(100).Render(view))
	return doc.String()
}

func (v *GameView) Enter() tea.Msg {
	switch v.GameState {
	case model.CutState:
		state.CutIndex = v.CutInput.Value()
		return services.CutDeck
	case model.DiscardState:
		p, err := utils.GetPlayerId(state.AccountId, state.ActiveMatch.Players)

		if err != nil {
			utils.Logger.Sugar().Error(err)
		}
		state.CurrentHandModifier = model.HandModifier{
			MatchId:  state.ActiveMatchId,
			PlayerId: p.Id,
			CardIds:  v.HighlightedIds,
		}
		return services.PutKitty
	}

	return nil
}

func (v *GameView) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	v.CutInput.Focus()
	v.CutInput, cmd = v.CutInput.Update(msg)
	return cmd
}

func (v *GameView) UpdateState(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	if v.PreviousState == v.GameState {
		return nil
	}

	v.PreviousState = v.GameState

	matchMsg := services.GetPlayerMatch()
	var match *model.GameMatch
	if err := json.Unmarshal(matchMsg.([]byte), &match); err != nil {
		return nil
	}

	p := utils.GetPlayerForAccountId(state.AccountId, match)

	for _, cardId := range p.Hand {
		card := utils.GetCardById(cardId)
		if card != nil {
			v.Hand = append(v.Hand, *card)
		}
	}

	for _, cardId := range p.Kitty {
		card := utils.GetCardById(cardId)
		v.Kitty = append(v.Kitty, *card)
	}

	for _, cardId := range p.Play {
		card := utils.GetCardById(cardId)
		v.Play = append(v.Play, *card)
	}

	return cmd
}
