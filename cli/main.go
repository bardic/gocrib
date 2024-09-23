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
				for _, idx := range m.SelectedCardIds {
					m.Kitty = append(m.Kitty, getCardById(idx, m.Hand))
					m.Hand = slices.DeleteFunc(m.Hand, func(c model.Card) bool {
						return c.Id == idx
					})
				}
			} else {
				for _, idx := range m.SelectedCardIds {
					m.CardsInPlay = append(m.CardsInPlay, getCardById(idx, m.Hand))
					m.Hand = slices.DeleteFunc(m.Hand, func(c model.Card) bool {
						return c.Id == idx
					})
				}
			}
			return m, nil
		case " ":
			cards := m.Hand
			idx := slices.Index(m.SelectedCardIds, cards[m.SelectedCardId].Id)
			if idx > -1 {
				m.SelectedCardIds = slices.Delete(m.SelectedCardIds, idx, 1)
			} else {
				m.SelectedCardIds = append(m.SelectedCardIds, cards[m.SelectedCardId].Id)
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
				m.SelectedCardId = 1
				m.ActiveSlot = model.CardTwo
			case model.CardTwo:
				m.SelectedCardId = 2
				m.ActiveSlot = model.CardThree
			case model.CardThree:
				m.SelectedCardId = 3
				m.ActiveSlot = model.CardFour
			case model.CardFour:
				if len(m.Hand) == 4 {
					m.ActiveSlot = model.CardOne
					m.SelectedCardId = 0
				} else {
					m.ActiveSlot = model.CardFive
					m.SelectedCardId = 4
				}
			case model.CardFive:
				m.SelectedCardId = 5
				m.ActiveSlot = model.CardSix
			case model.CardSix:
				m.SelectedCardId = 0
				m.ActiveSlot = model.CardOne
			}
		case "left":
			switch m.ActiveSlot {
			case model.CardOne:
				if len(m.Hand) == 4 {
					m.SelectedCardId = 3
					m.ActiveSlot = model.CardFour
				} else {
					m.SelectedCardId = 5
					m.ActiveSlot = model.CardSix
				}
			case model.CardTwo:
				m.SelectedCardId = 0
				m.ActiveSlot = model.CardOne
			case model.CardThree:
				m.SelectedCardId = 1
				m.ActiveSlot = model.CardTwo
			case model.CardFour:
				m.SelectedCardId = 2
				m.ActiveSlot = model.CardThree
			case model.CardFive:
				m.SelectedCardId = 3
				m.ActiveSlot = model.CardFour
			case model.CardSix:
				m.SelectedCardId = 4
				m.ActiveSlot = model.CardFive
			}
		}
	case timer.TickMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		state.ActiveMatchId = 2
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

		diffs := match.Eq(state.ActiveMatch)
		if diffs == 0 {
			return m, tea.Batch(cmds...)
		}

		for diff := model.GenericDiff; diff < model.MaxDiff; diff <<= 1 {
			d := diffs & diff
			if d != 0 {
				switch d {
				case model.CutDiff:
					m.ViewState = model.BoardView
				case model.CardsInPlayDiff:

				case model.GameStateDiff:
					//cmds = append(cmds, updateView)
				case model.GenericDiff:
				case model.NewDeckDiff:
					m.ViewState = model.BoardView
				case model.MaxDiff:
				case model.TurnDiff:
				case model.TurnPassTimestampsDiff:
				}
			}
		}

		state.ActiveMatch = match
		state.ActiveViewModel = m.ViewModel
		//Once it does change, determine the change

	case views.ViewModel:
		switch msg.ViewState {
		case model.LobbyView:
		case model.PlayView:
		case model.KittyView:
		case model.HandView:
		case model.ScoresView:
		case model.GameOverView:
		}
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
		v = styles.ViewStyle.Render(views.GameView(m.ViewModel))
	case model.PlayView:
		v = styles.ViewStyle.Render(views.GameView(m.ViewModel))
	case model.HandView:
		v = styles.ViewStyle.Render(views.GameView(m.ViewModel))
	case model.KittyView:
		return styles.ViewStyle.Render(views.GameView(m.ViewModel))
	default:
		v = styles.ViewStyle.Render(views.LobbyView())
	}
	return v
}

func newModel() appModel {
	m := appModel{
		ViewModel: views.ViewModel{
			ActiveSlot:  model.CardOne,
			ViewState:   model.ActiveView,
			CardsInPlay: []model.Card{},
			Hand: []model.Card{
				{Id: 1, Suit: 0, Value: 1, Art: "meow.png"},
				{Id: 2, Suit: 1, Value: 2, Art: "meow.png"},
				{Id: 3, Suit: 2, Value: 3, Art: "meow.png"},
				{Id: 4, Suit: 3, Value: 4, Art: "meow.png"},
				{Id: 5, Suit: 3, Value: 5, Art: "meow.png"},
				{Id: 6, Suit: 3, Value: 6, Art: "meow.png"},
			},
			SelectedCardId: 0,
			Tabs:           model.TabNames,
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

	return match
}

func getCardById(id int, hand []model.Card) model.Card {
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
