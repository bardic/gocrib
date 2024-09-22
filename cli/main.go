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
	timer timer.Model
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
			cmds = append(cmds, updateView)
		case " ":
			cards := m.Hand
			if m.ActiveTab == 2 {
				cards = m.Kitty
			}
			idx := slices.Index(m.SelectedCardIds, cards[m.Next].Id)
			if idx > -1 {
				m.SelectedCardIds = slices.Delete(m.SelectedCardIds, idx, 1)
			} else {
				m.SelectedCardIds = append(m.SelectedCardIds, cards[m.Next].Id)
			}
		case "tab":
			m.SelectedCardIds = make([]int, 0)
			m.ActiveTab = min(m.ActiveTab+1, len(m.Tabs)-1)
			return m, nil
		case "shift+tab":
			m.SelectedCardIds = make([]int, 0)
			m.ActiveTab = max(m.ActiveTab-1, 0)
			return m, nil
		case "right":
			switch m.ActiveSlot {
			case model.CardOne:
				m.Next = 1
				m.ActiveSlot = model.CardTwo
			case model.CardTwo:
				m.Next = 2
				m.ActiveSlot = model.CardThree
			case model.CardThree:
				m.Next = 3
				m.ActiveSlot = model.CardFour
			case model.CardFour:
				if len(m.Hand) == 4 {
					m.ActiveSlot = model.CardOne
					m.Next = 0
				} else {
					m.ActiveSlot = model.CardFive
					m.Next = 4
				}
			case model.CardFive:
				m.Next = 5
				m.ActiveSlot = model.CardSix
			case model.CardSix:
				m.Next = 0
				m.ActiveSlot = model.CardOne
			}
		case "left":
			switch m.ActiveSlot {
			case model.CardOne:
				if len(m.Hand) == 4 {
					m.Next = 3
					m.ActiveSlot = model.CardFour
				} else {
					m.Next = 5
					m.ActiveSlot = model.CardSix
				}
			case model.CardTwo:
				m.Next = 0
				m.ActiveSlot = model.CardOne
			case model.CardThree:
				m.Next = 1
				m.ActiveSlot = model.CardTwo
			case model.CardFour:
				m.Next = 2
				m.ActiveSlot = model.CardThree
			case model.CardFive:
				m.Next = 3
				m.ActiveSlot = model.CardFour
			case model.CardSix:
				m.Next = 4
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
					// fmt.Println("CutDiff")
					m.GameState = model.CutState
					m.ViewState = model.CutView
					cmds = append(cmds, updateView)
				case model.CardsInPlayDiff:
					// fmt.Println("CardsInPlayDiff")
					if m.GameState&model.OpponentState == 1 {
						m.GameState = model.PlayState
					} else {
						m.GameState = model.OpponentState
					}

					cmds = append(cmds, updateView)
				case model.GameStateDiff:
					// fmt.Println("GameStateDiff" + string(m.GameState))
					cmds = append(cmds, updateView)
				case model.GenericDiff:
					// fmt.Println("GenericDiff")

				case model.MaxDiff:
					// fmt.Println("MaxDiff")
				case model.TurnDiff:
					// fmt.Println("TurnDiff")
				case model.TurnPassTimestampsDiff:
					// fmt.Println("TurnPassTimestampsDiff")
				}
			}
		}

		state.ActiveMatch = match
		//Once it does change, determine the change

	case appModel:
		switch msg.ViewState {
		case model.LobbyView:
			m.ViewState = model.LobbyView
			m.ActiveTab = 1
		case model.PlayView:
			m.GameState = model.PlayState
			m.ActiveTab = 2
		case model.KittyView:
			m.GameState = model.PlayState
			m.ActiveTab = 3
		case model.HandView:
			m.GameState = model.PlayState
			m.ActiveTab = 1
		case model.ScoresView:
			break
		case model.GameOverView:
			break
		}
	}
	return m, tea.Batch(cmds...)
}

func updateView() tea.Msg {
	//fmt.Println("update view")
	m := state.ActiveViewModel
	if m.ViewState == model.PlayView {
		switch m.GameState {
		case model.DealState:
			break
		case model.CutState:

			break
		case model.DiscardState:
			for _, idx := range m.SelectedCardIds {
				m.Kitty = append(m.Kitty, getCardById(idx, m.Hand))
				m.Hand = slices.DeleteFunc(m.Hand, func(c model.Card) bool {
					return c.Id == idx
				})
			}
			m.GameState = model.PlayState
			m.ActiveTab = 2
		case model.PlayState:
			for _, idx := range m.SelectedCardIds {
				m.CardsInPlay = append(m.CardsInPlay, getCardById(idx, m.Hand))
				m.Hand = slices.DeleteFunc(m.Hand, func(c model.Card) bool {
					return c.Id == idx
				})
			}
			m.GameState = model.OpponentState
			m.ActiveTab = 0
		case model.OpponentState:
			break
		case model.KittyState:
			break
		case model.GameWonState:
			break
		case model.GameLostState:
			break
		}
		m.SelectedCardIds = make([]int, 0)
	} else {
		m.ViewState = model.PlayView
	}

	return ""
}

func (m appModel) View() string {
	v := styles.ViewStyle.Render(views.LobbyView())

	switch m.ViewState {
	case model.LobbyView:
		v = styles.ViewStyle.Render(views.LobbyView())
	case model.PlayView:
		v = styles.ViewStyle.Render(views.GameView(m.ViewModel))
	case model.KittyView:
		//return styles.ViewStyle.Render(m.kittyView())
		break
	case model.CutView:
		v = styles.ViewStyle.Render(views.CutView())
	case model.HandView:
		//return styles.ViewStyle.Render(m.HandView())
		break
	default:
		v = styles.ViewStyle.Render(views.LobbyView())
	}
	return v
}

func newModel() appModel {
	m := appModel{
		ViewModel: views.ViewModel{
			ActiveSlot:  model.CardOne,
			ViewState:   model.LobbyView,
			GameState:   model.WaitingState,
			CardsInPlay: []model.Card{},
			Hand: []model.Card{
				{Id: 1, Suit: 0, Value: 1, Art: "meow.png"},
				{Id: 2, Suit: 1, Value: 2, Art: "meow.png"},
				{Id: 3, Suit: 2, Value: 3, Art: "meow.png"},
				{Id: 4, Suit: 3, Value: 4, Art: "meow.png"},
				{Id: 5, Suit: 3, Value: 5, Art: "meow.png"},
				{Id: 6, Suit: 3, Value: 6, Art: "meow.png"},
			},
			Next: 0,
			Tabs: []string{"Play", "Hand", "Kitty"},
		},

		timer: timer.NewWithInterval(time.Hour, time.Second),
	}
	return m
}

func createMatch() model.Match {
	newMatch := services.PostPlayerMatch().([]byte)

	var match model.Match
	json.Unmarshal(newMatch, &match)

	return match
}

func createGame() tea.Msg {
	match := createMatch()

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
