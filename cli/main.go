package main

import (
	"encoding/json"
	"fmt"
	"log"
	"slices"
	"strconv"
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
	hand      []model.Card
	kitty     []model.Card
	play      []model.Card
	gameState model.GameState
	timer     timer.Model
}

func (m appModel) Init() tea.Cmd {
	return m.timer.Init()
}

func (m appModel) OnEnterDuringPlay() (appModel, tea.Cmd) {
	if m.gameState == model.WaitingState {
		m.gameState = model.DiscardState
	}

	m.ViewState = model.GameView
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

		m.HighlightedIds = []int{}
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

		m.HighlightedIds = []int{}
		return m, services.PutPlay
	}
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
			switch m.ViewState {
			case model.LoginView:
				idStr := views.LoginIdField.Value()
				id, err := strconv.Atoi(idStr)
				if err != nil {
					return m, tea.Quit
				}

				state.PlayerId = id
				return m, services.Login
			case model.LobbyView:
				idStr := views.LobbyTable.SelectedRow()[0]
				id, err := strconv.Atoi(idStr)
				if err != nil {
					return m, tea.Quit
				}

				state.ActiveMatchId = id
				return m, services.JoinMatch
			case model.GameView:
				return m.OnEnterDuringPlay()
			case model.GameOverView:
			}
		case "n":
			createGame()
			m.ViewState = model.GameView
		case " ":
			cards := m.hand
			idx := slices.Index(m.HighlightedIds, cards[m.HighlighedId].Id)
			if idx > -1 {
				m.HighlightedIds = slices.Delete(m.HighlightedIds, idx, 1)
			} else {
				m.HighlightedIds = append(m.HighlightedIds, cards[m.HighlighedId].Id)
			}
		case "tab":
			switch m.ViewState {
			case model.LobbyView:
				m.ActiveLandingTab = m.ActiveLandingTab + 1
				switch m.ActiveLandingTab {
				case 0:
					m.LobbyViewState = model.OpenMatches
				case 1:
					m.LobbyViewState = model.AvailableMatches

				}
			case model.GameView:
				m.ActiveTab = m.ActiveTab + 1
				switch m.ActiveTab {
				case 0:
					m.GameViewState = model.BoardView
				case 1:
					m.GameViewState = model.PlayView
				case 2:
					m.GameViewState = model.HandView
				case 3:
					m.GameViewState = model.KittyView
				}
			}
		case "shift+tab":
			switch m.ViewState {
			case model.LobbyView:
				m.ActiveLandingTab = m.ActiveLandingTab - 1
				switch m.ActiveLandingTab {
				case 0:
					m.LobbyViewState = model.OpenMatches
				case 1:
					m.LobbyViewState = model.AvailableMatches

				}
			case model.GameView:
				m.ActiveTab = m.ActiveTab - 1
				switch m.ActiveTab {
				case 0:
					m.GameViewState = model.BoardView
				case 1:
					m.GameViewState = model.PlayView
				case 2:
					m.GameViewState = model.HandView
				case 3:
					m.GameViewState = model.KittyView
				}
			}
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
		var cmd tea.Cmd

		if m.ViewState == model.LoginView {
			views.LoginIdField.Focus()
			views.LoginIdField, cmd = views.LoginIdField.Update(msg)
		}

		if m.ViewState == model.LobbyView {
			views.LobbyTable.Focus()
			views.LobbyTable, cmd = views.LobbyTable.Update(msg)

		}

		cmds = append(cmds, cmd)
	case timer.TickMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)

		if m.ViewState == model.LobbyView || m.ViewState == model.LoginView {
			cmds = append(cmds, cmd)
		} else {
			cmds = append(cmds, cmd, services.GetPlayerMatch)
		}
	case []byte:
		var matchStr string

		err := json.Unmarshal(msg, &matchStr)

		if err != nil {
			fmt.Println(err)
		}

		switch m.ViewState {
		case model.LoginView:
			m.ViewState = model.LobbyView
			var account model.Account
			err = json.Unmarshal([]byte(matchStr), &account)

			if err != nil {
				fmt.Println(err)
			}

			state.PlayerId = account.Id

			return m, tea.Batch(cmds...)
		case model.LobbyView:
			m.ViewState = model.GameView
			return m, tea.Batch(cmds...)
		}

		var match model.GameMatch
		err = json.Unmarshal([]byte(matchStr), &match)

		if err != nil {
			fmt.Println(err)
		}

		// if match.GameState == m.gameState {
		// 	return m, nil
		// }

		diffs := match.Eq(state.ActiveMatch.Match)

		if diffs == 0 {
			return m, tea.Batch(cmds...)
		}

		state.ActiveMatch = match

		m.hand = []model.Card{}
		m.kitty = []model.Card{}
		m.play = []model.Card{}

		for _, cardId := range match.Players[0].Hand {
			card := getCardById(cardId)
			m.hand = append(m.hand, card)
		}

		for _, cardId := range match.Players[0].Kitty {
			card := getCardById(cardId)
			m.kitty = append(m.hand, card)
		}

		for _, cardId := range match.Players[0].Play {
			card := getCardById(cardId)
			m.play = append(m.hand, card)
		}

		for diff := model.GenericDiff; diff < model.MaxDiff; diff <<= 1 {
			d := diffs & diff
			if d != 0 {
				switch d {
				case model.CutDiff:
					fmt.Println("cutdiff")
					m.ViewState = model.GameView
				case model.CardsInPlayDiff:
					fmt.Println("cards in play diff")
				case model.GameStateDiff:
					fmt.Println("game state diff")
				case model.GenericDiff:
					fmt.Println("generic diff")
				case model.NewDeckDiff:
					m.ViewState = model.GameView
					fmt.Println("new deck diff")
				case model.MaxDiff:
					fmt.Println("max diff")
				case model.TurnDiff:
					fmt.Println("turn diff")
				case model.TurnPassTimestampsDiff:
					fmt.Println("pass timestamp diff")
				}
			}
		}
		m.gameState = match.GameState
		state.ActiveMatch = match
	case model.GameMatch:
		deckByte := services.GetDeckById(msg.DeckId).([]byte)
		var deckJson string
		json.Unmarshal(deckByte, &deckJson)

		var deck model.GameDeck
		json.Unmarshal([]byte(deckJson), &deck)

		state.ActiveDeck = deck
		state.ActiveMatch = msg
	case model.Account:
		fmt.Println("cat")
	}

	return m, tea.Batch(cmds...)
}

func (m appModel) View() string {
	var v string
	switch m.ViewState {
	case model.LoginView:
		v = styles.ViewStyle.Render(views.LoginView())
	case model.LobbyView:
		view, err := views.LobbyView(m.ViewModel)
		if err != nil {
			return err.Error()
		}
		v = styles.ViewStyle.Render(view)
	case model.GameView:
		switch m.GameViewState {
		case model.BoardView:
			v = styles.ViewStyle.Render(views.GameView(
				m.HighlighedId,
				m.HighlightedIds,
				[]model.Card{},
				m.ViewModel,
				m.gameState))
		case model.PlayView:
			v = styles.ViewStyle.Render(views.GameView(
				m.HighlighedId,
				m.HighlightedIds,
				m.play,
				m.ViewModel,
				m.gameState))
		case model.HandView:
			v = styles.ViewStyle.Render(views.GameView(m.HighlighedId,
				m.HighlightedIds,
				m.hand,
				m.ViewModel,
				m.gameState))
		case model.KittyView:
			v = styles.ViewStyle.Render(views.GameView(m.HighlighedId,
				m.HighlightedIds,
				m.kitty,
				m.ViewModel,
				m.gameState))
		}
	default:
		view, err := views.LobbyView(m.ViewModel)
		if err != nil {
			return err.Error()
		}
		//m.table = table
		v = styles.ViewStyle.Render(view)
	}
	return v
}

func newModel() appModel {
	m := appModel{
		ViewModel: views.ViewModel{
			ActiveSlot:     model.CardOne,
			ViewState:      model.LoginView,
			Tabs:           model.GameTabNames,
			LobbyTabs:      model.LobbyTabNames,
			HighlighedId:   0,
			HighlightedIds: []int{},
		},
		hand:      []model.Card{},
		gameState: model.WaitingState,
		timer:     timer.NewWithInterval(time.Hour, time.Second*1),
	}
	return m
}

func createGame() tea.Msg {
	newMatch := services.PostPlayerMatch().([]byte)

	var match model.GameMatch
	json.Unmarshal(newMatch, &match)
	state.ActiveMatchId = match.Id
	state.ActiveMatch = match
	return match
}

func getCardById(id int) model.Card {
	for _, card := range state.ActiveDeck.Cards {
		if card.Id == id {
			return card.Card
		}
	}
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
