package main

import (
	"encoding/json"
	"fmt"
	"log"
	"slices"
	"time"

	"github.com/bardic/cribbagev2/cli/services"
	"github.com/bardic/cribbagev2/cli/state"
	"github.com/bardic/cribbagev2/cli/styles"
	"github.com/bardic/cribbagev2/cli/views"
	"github.com/bardic/cribbagev2/model"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
)

type appModel struct {
	views.ViewModel
	hand  []model.Card
	kitty []model.Card
	play  []model.Card

	gameState model.GameState
	timer     timer.Model
}

func (m appModel) Init() tea.Cmd {
	return tea.Batch(createGame, m.timer.Init())
}

func (m appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "c":
			return m, tea.Quit
		case "enter", "view_update":
			if m.gameState == model.DiscardState {
				for _, idx := range m.HighlightedIds {
					m.kitty = append(m.kitty, getCardInHandById(idx, m.hand))
					m.hand = slices.DeleteFunc(m.hand, func(c model.Card) bool {
						return c.Id == idx
					})
				}

				state.CurrentHandModifier = model.HandModifier{
					MatchId:  state.ActiveMatchId,
					CardIds:  getIdsFromCards(m.kitty),
					PlayerId: state.ActiveMatch.PlayerIds[0],
				}

				return m, services.PutKitty
			} else {
				for _, idx := range m.HighlightedIds {
					m.play = append(m.play, getCardInHandById(idx, m.hand))
					m.hand = slices.DeleteFunc(m.hand, func(c model.Card) bool {
						return c.Id == idx
					})
				}

				state.CurrentHandModifier = model.HandModifier{
					MatchId:  state.ActiveMatchId,
					CardIds:  getIdsFromCards(m.play),
					PlayerId: state.ActiveMatch.PlayerIds[0],
				}

				return m, services.PutPlay
			}
		case " ":
			cards := m.hand
			idx := slices.Index(m.HighlightedIds, cards[m.HighlighedId].Id)
			if idx > -1 {
				m.HighlightedIds = slices.Delete(m.HighlightedIds, idx, 1)
			} else {
				m.HighlightedIds = append(m.HighlightedIds, cards[m.HighlighedId].Id)
			}
		case "tab":
			m.ActiveTab = m.ActiveTab + 1
			switch m.ActiveTab {
			case 0:
				m.ViewState = model.BoardView
			case 1:
				m.ViewState = model.PlayView
			case 2:
				m.ViewState = model.HandView
			case 3:
				m.ViewState = model.KittyView
			}
			return m, nil
		case "shift+tab":
			m.ActiveTab = m.ActiveTab - 1
			switch m.ActiveTab {
			case 0:
				m.ViewState = model.BoardView
			case 1:
				m.ViewState = model.PlayView
			case 2:
				m.ViewState = model.HandView
			case 3:
				m.ViewState = model.KittyView
			}
			return m, nil
		case "right":
			switch m.ActiveSlot {
			case model.CardOne:
				m.HighlighedId = 1
				m.ActiveSlot = model.CardTwo
			case model.CardTwo:
				m.HighlighedId = 2
				m.ActiveSlot = model.CardThree
			case model.CardThree:
				m.HighlighedId = 3
				m.ActiveSlot = model.CardFour
			case model.CardFour:
				if len(m.hand) == 4 {
					m.ActiveSlot = model.CardOne
					m.HighlighedId = 0
				} else {
					m.ActiveSlot = model.CardFive
					m.HighlighedId = 4
				}
			case model.CardFive:
				m.HighlighedId = 5
				m.ActiveSlot = model.CardSix
			case model.CardSix:
				m.HighlighedId = 0
				m.ActiveSlot = model.CardOne
			}
		case "left":
			switch m.ActiveSlot {
			case model.CardOne:
				if len(m.hand) == 4 {
					m.HighlighedId = 3
					m.ActiveSlot = model.CardFour
				} else {
					m.HighlighedId = 5
					m.ActiveSlot = model.CardSix
				}
			case model.CardTwo:
				m.HighlighedId = 0
				m.ActiveSlot = model.CardOne
			case model.CardThree:
				m.HighlighedId = 1
				m.ActiveSlot = model.CardTwo
			case model.CardFour:
				m.HighlighedId = 2
				m.ActiveSlot = model.CardThree
			case model.CardFive:
				m.HighlighedId = 3
				m.ActiveSlot = model.CardFour
			case model.CardSix:
				m.HighlighedId = 4
				m.ActiveSlot = model.CardFive
			}
		}
	case timer.TickMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)

		cmds = append(cmds, cmd, services.GetPlayerMatch)
	case []byte:
		var matchStr string

		err := json.Unmarshal(msg, &matchStr)

		if err != nil {
			fmt.Println(err)
		}

		var match model.Match
		err = json.Unmarshal([]byte(matchStr), &match)

		if err != nil {
			fmt.Println(err)
		}

		if match.GameState == m.gameState {
			return m, nil
		}

		diffs := match.Eq(state.ActiveMatch)

		if diffs == 0 {
			return m, tea.Batch(cmds...)
		}

		for diff := model.GenericDiff; diff < model.MaxDiff; diff <<= 1 {
			d := diffs & diff
			if d != 0 {
				switch d {
				case model.CutDiff:
					// fmt.Println("cutdiff")
					m.ViewState = model.BoardView
				case model.CardsInPlayDiff:
					// fmt.Println("cards in play diff")
				case model.GameStateDiff:
					// fmt.Println("game state diff")
				case model.GenericDiff:
					// fmt.Println("generic diff")
				case model.NewDeckDiff:
					m.ViewState = model.BoardView
					// fmt.Println("new deck diff")
				case model.MaxDiff:
					// fmt.Println("max diff")
				case model.TurnDiff:
					// fmt.Println("turn diff")
				case model.TurnPassTimestampsDiff:
					// fmt.Println("pass timestamp diff")
				}
			}
		}
		m.gameState = match.GameState
	case model.Match:
		deckByte := services.GetDeckById(msg.DeckId).([]byte)
		var deckJson string
		json.Unmarshal(deckByte, &deckJson)

		var deck model.GameDeck
		json.Unmarshal([]byte(deckJson), &deck)

		state.ActiveDeck = deck
		//	state.ActiveMatchId = msg.Id

		m.gameState = msg.GameState
		m.ViewState = model.BoardView
	}
	return m, tea.Batch(cmds...)
}

func (m appModel) View() string {
	var v string
	switch m.ViewState {
	case model.ActiveView:
		v = styles.ViewStyle.Render(views.ActiveView())
	case model.LobbyView:
		v = styles.ViewStyle.Render(views.LobbyView())
	case model.BoardView:
		v = styles.ViewStyle.Render(views.GameView(
			m.HighlighedId,
			m.HighlightedIds,
			[]model.Card{},
			m.ViewModel))
	case model.PlayView:
		v = styles.ViewStyle.Render(views.GameView(m.HighlighedId,
			m.HighlightedIds,
			m.play,
			m.ViewModel))
	case model.HandView:
		v = styles.ViewStyle.Render(views.GameView(m.HighlighedId,
			m.HighlightedIds,
			m.hand,
			m.ViewModel))
	case model.KittyView:
		return styles.ViewStyle.Render(views.GameView(m.HighlighedId,
			m.HighlightedIds,
			m.kitty,
			m.ViewModel))
	default:
		v = styles.ViewStyle.Render(views.LobbyView())
	}
	return v
}

func newModel() appModel {
	m := appModel{
		ViewModel: views.ViewModel{
			ActiveSlot:     model.CardOne,
			ViewState:      model.ActiveView,
			Tabs:           model.TabNames,
			HighlighedId:   0,
			HighlightedIds: []int{},
		},
		hand: []model.Card{
			{Id: 1, Suit: 0, Value: 1, Art: "meow.png"},
			{Id: 2, Suit: 1, Value: 2, Art: "meow.png"},
			{Id: 3, Suit: 2, Value: 3, Art: "meow.png"},
			{Id: 4, Suit: 3, Value: 4, Art: "meow.png"},
			{Id: 5, Suit: 3, Value: 5, Art: "meow.png"},
			{Id: 6, Suit: 3, Value: 6, Art: "meow.png"},
		},
		gameState: model.WaitingState,
		timer:     timer.NewWithInterval(time.Hour, time.Second*1),
	}
	return m
}

func createGame() tea.Msg {
	newMatch := services.PostPlayerMatch().([]byte)

	var match model.Match
	json.Unmarshal(newMatch, &match)
	state.ActiveMatchId = match.Id
	state.ActiveMatch = match
	return match
}

func getCardById(id int) model.Card {
	return model.Card{}
}

func getIdsFromCards(c []model.Card) []int {
	ids := make([]int, len(c))
	for _, card := range c {
		ids = append(ids, card.Id)
	}

	return ids
}

func getCardInHandById(id int, hand []model.Card) model.Card {
	idx := slices.IndexFunc(hand, func(c model.Card) bool {
		return c.Id == id
	})

	if idx == -1 {
		return model.Card{}
	}

	return hand[idx]
}

func main() {
	p := tea.NewProgram(newModel())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
