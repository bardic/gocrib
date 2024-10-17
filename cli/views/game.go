package views

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/bardic/gocrib/cli/services"
	"github.com/bardic/gocrib/cli/styles"
	"github.com/bardic/gocrib/cli/utils"
	"github.com/bardic/gocrib/model"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type GameView struct {
	Initd          bool
	MatchId        int
	AccountId      int
	CutIndex       string
	HighlightedIds []int
	GameTabNames   []string
	GameViewState  model.GameViewState
	ActiveSlot     int
	ActiveTab      int
	HighlighedId   int
	CutInput       textinput.Model
	Deck           *model.GameDeck
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
	if v.Initd {
		return
	}

	matchMsg := services.GetPlayerMatch(strconv.Itoa(v.MatchId))
	var match *model.GameMatch
	if err := json.Unmarshal(matchMsg.([]byte), &match); err != nil {
		return
	}

	v.Initd = true
	v.GameTabNames = []string{"Board", "Play", "Hand", "Kitty"}
	v.CutInput = textinput.New()
	v.CutInput.Placeholder = "0"
	v.CutInput.CharLimit = 5
	v.CutInput.Width = 5

	deckByte := services.GetDeckById(match.DeckId).([]byte)
	var deck model.GameDeck
	json.Unmarshal(deckByte, &deck)
	v.Deck = &deck
}

func (v *GameView) View() string {
	if v.GameMatch == nil {
		return "Loading..."
	}

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
		if v.GameMatch.GameState == model.CutState {
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
		p := utils.GetPlayerForAccountId(v.AccountId, v.GameMatch)
		view = HandView(v.HighlighedId, v.HighlightedIds, p.Play, v.Deck)
	case model.HandView:
		p := utils.GetPlayerForAccountId(v.AccountId, v.GameMatch)
		if v.GameMatch.GameState == 0 {
			view = "Waiting to be dealt"
		} else {
			view = HandView(v.HighlighedId, v.HighlightedIds, p.Hand, v.Deck)
		}
	case model.KittyView:
		p := utils.GetPlayerForAccountId(v.AccountId, v.GameMatch)
		if len(p.Kitty) == 0 {
			view = "Empty Kitty"
		} else {
			view = HandView(v.HighlighedId, v.HighlightedIds, p.Kitty, v.Deck)
		}
	}

	doc.WriteString(styles.WindowStyle.Width(100).Render(view))
	return doc.String()
}

func (v *GameView) Enter() tea.Msg {
	switch v.GameMatch.GameState {
	case model.CutState:
		v.CutIndex = v.CutInput.Value()
		return services.CutDeck
	case model.DiscardState:
		p, err := utils.GetPlayerId(v.AccountId, v.GameMatch.Players)

		if err != nil {
			utils.Logger.Sugar().Error(err)
		}

		services.PutKitty(model.HandModifier{
			MatchId:  v.GameMatch.Id,
			PlayerId: p.Id,
			CardIds:  v.HighlightedIds,
		})
	}

	return nil
}

func (v *GameView) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	v.CutInput.Focus()
	v.CutInput, cmd = v.CutInput.Update(msg)
	return cmd
}

/*func (v *GameView) UpdateState(newState model.GameState) tea.Cmd {
	var cmd tea.Cmd

	if v.GameMatch == nil || v.GameMatch.GameState == newState {
		return nil
	}

	matchMsg := services.GetPlayerMatch(strconv.Itoa(v.MatchId))
	var match *model.GameMatch
	if err := json.Unmarshal(matchMsg.([]byte), &match); err != nil {
		return nil
	}

	v.GameMatch = match
	p := utils.GetPlayerForAccountId(v.AccountId, match)

	for _, cardId := range p.Hand {
		card := utils.GetCardById(cardId, v.Deck)
		if card != nil {
			v.Hand = append(v.Hand, *card)
		}
	}

	for _, cardId := range p.Kitty {
		card := utils.GetCardById(cardId, v.Deck)
		v.Kitty = append(v.Kitty, *card)
	}

	for _, cardId := range p.Play {
		card := utils.GetCardById(cardId, v.Deck)
		v.Play = append(v.Play, *card)
	}

	return cmd
}*/
