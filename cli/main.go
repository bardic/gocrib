package main

import (
	"encoding/json"
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/bardic/cribbagev2/cli/services"
	"github.com/bardic/cribbagev2/model"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type cardSlots uint

const (
	cardOne cardSlots = iota
	cardTwo
	cardThree
	cardFour
	cardFive
	cardSix
)

var (
	modelStyle = lipgloss.NewStyle().
			Width(10).
			Height(5).
			Align(lipgloss.Center, lipgloss.Center).
			BorderStyle(lipgloss.HiddenBorder())

	selectedStyle = lipgloss.NewStyle().
			Width(10).
			Height(5).
			Align(lipgloss.Center, lipgloss.Center).
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(lipgloss.Color("10"))

	selectedFocusedStyle = lipgloss.NewStyle().
				Width(10).
				Height(5).
				Align(lipgloss.Center, lipgloss.Center).
				BorderStyle(lipgloss.ThickBorder()).
				BorderForeground(lipgloss.Color("69"))

	focusedModelStyle = lipgloss.NewStyle().
				Width(10).
				Height(5).
				Align(lipgloss.Center, lipgloss.Center).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("69"))

	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")

	viewStyle        = lipgloss.NewStyle().Padding(1, 2, 1, 2)
	highlightColor   = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 1)
	activeTabStyle   = inactiveTabStyle.Border(activeTabBorder, true)
	windowStyle      = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(2, 0).Border(lipgloss.NormalBorder()).UnsetBorderTop()
)

type mainModel struct {
	activeSlot      cardSlots
	gameState       model.GameState
	viewState       model.ViewState
	hand            []model.Card
	kitty           []model.Card
	cardsInPlay     []model.Card
	selectedCardIds []int
	next            int
	activeTab       int
	tabs            []string
}

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

func newModel() mainModel {
	m := mainModel{
		activeSlot:  cardOne,
		viewState:   model.LobbyView,
		gameState:   model.WaitingState,
		cardsInPlay: []model.Card{},
		hand: []model.Card{
			{Id: 1, Suit: 0, Value: 1, Art: "meow.png"},
			{Id: 2, Suit: 1, Value: 2, Art: "meow.png"},
			{Id: 3, Suit: 2, Value: 3, Art: "meow.png"},
			{Id: 4, Suit: 3, Value: 4, Art: "meow.png"},
			{Id: 5, Suit: 3, Value: 5, Art: "meow.png"},
			{Id: 6, Suit: 3, Value: 6, Art: "meow.png"},
		},
		next: 0,
		tabs: []string{"Play", "Hand", "Kitty"},
	}
	return m
}

// func createDeck() model.GameDeck {
// 	deck := model.GameDeck{
// 		Cards: []int{
// 			1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52,
// 		},
// 	}

// 	newDeck := services.PostDeck(deck).([]byte)
// 	var deckId int
// 	json.Unmarshal(newDeck, &deckId)
// 	deck.Id = deckId

// 	return deck
// }

func createMatch() model.Match {

	newMatch := services.PostPlayerMatch().([]byte)

	var match model.Match
	json.Unmarshal(newMatch, &match)

	return match
}

func createPlayerForMatch(matchId int, accountId int) model.Player {
	player := model.Player{
		AccountId: accountId,
	}

	newPlayer := services.PostPlayer(player).([]byte)

	fmt.Println(string(newPlayer))

	var playerId int
	json.Unmarshal(newPlayer, &playerId)
	player.Id = matchId

	return player
}

func createGame() tea.Msg {
	// deck := createDeck()
	createMatch()
	match := createMatch()

	fmt.Println(match)

	player1 := createPlayerForMatch(match.Id, 1)
	player2 := createPlayerForMatch(match.Id, 2)

	fmt.Println(player1)

	match.PlayerIds = []int{player1.Id, player2.Id}

	return match
}

func (m mainModel) Init() tea.Cmd {
	return createGame

	// return createMatch

	// //create match

	// //Start Match
	// // when second player joins server will dispatch a state change message
	// //set game state from waiting to discard

	// m.gameState = model.GameState(model.Discard)

	// submit dicarded cards

	// enter free play state

	// Check for round end
	// Check if freeplay end
	// reset player turn and start next round

	// count hand
	// loop

	// //get all cards

	// // r = services.GetAllCards()
	// // var c string

	// // if err := json.Unmarshal((r.([]byte)), &c); err != nil {
	// // 	panic(err)
	// // }

	// // var cards []model.Card

	// // if err := json.Unmarshal(([]byte(c)), &cards); err != nil {
	// // 	panic(err)
	// // }

	// //start match

	// // submit action

	// a := model.GameAction{
	// 	MatchId:        match.Id,
	// 	Type:           model.Discard,
	// 	GameplayCardId: cards[0].Id,
	// }
	// services.PostGame(a)

	// return nil
}

func handleErr(r tea.Msg) error {
	switch r := r.(type) {
	case error:
		return r
	}

	return nil
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if m.viewState == model.PlayView {
				switch m.gameState {
				case model.DealState:
					break
				case model.CutState:
					break
				case model.DiscardState:
					for _, idx := range m.selectedCardIds {
						m.kitty = append(m.kitty, m.getCardById(idx))
						m.hand = slices.DeleteFunc(m.hand, func(c model.Card) bool {
							return c.Id == idx
						})
					}
					m.gameState = model.PlayState
					m.activeTab = 2
				case model.PlayState:
					for _, idx := range m.selectedCardIds {
						m.cardsInPlay = append(m.cardsInPlay, m.getCardById(idx))
						m.hand = slices.DeleteFunc(m.hand, func(c model.Card) bool {
							return c.Id == idx
						})
					}
					m.gameState = model.OpponentState
					m.activeTab = 0
				case model.OpponentState:
					break
				case model.KittyState:
					break
				case model.GameWonState:
					break
				case model.GameLostState:
					break
				}
				m.selectedCardIds = make([]int, 0)
			} else {
				m.viewState = model.PlayView
			}
		case " ":
			cards := m.hand
			if m.activeTab == 2 {
				cards = m.kitty
			}
			idx := slices.Index(m.selectedCardIds, cards[m.next].Id)
			if idx > -1 {
				m.selectedCardIds = slices.Delete(m.selectedCardIds, idx, 1)
			} else {
				m.selectedCardIds = append(m.selectedCardIds, cards[m.next].Id)
			}
		case "tab":
			m.selectedCardIds = make([]int, 0)
			m.activeTab = min(m.activeTab+1, len(m.tabs)-1)
			return m, nil
		case "shift+tab":
			m.selectedCardIds = make([]int, 0)
			m.activeTab = max(m.activeTab-1, 0)
			return m, nil
		case "right":
			switch m.activeSlot {
			case cardOne:
				m.next = 1
				m.activeSlot = cardTwo
			case cardTwo:
				m.next = 2
				m.activeSlot = cardThree
			case cardThree:
				m.next = 3
				m.activeSlot = cardFour
			case cardFour:
				if len(m.hand) == 4 {
					m.activeSlot = cardOne
					m.next = 0
				} else {
					m.activeSlot = cardFive
					m.next = 4
				}
			case cardFive:
				m.next = 5
				m.activeSlot = cardSix
			case cardSix:
				m.next = 0
				m.activeSlot = cardOne
			}
		case "left":
			switch m.activeSlot {
			case cardOne:
				if len(m.hand) == 4 {
					m.next = 3
					m.activeSlot = cardFour
				} else {
					m.next = 5
					m.activeSlot = cardSix
				}
			case cardTwo:
				m.next = 0
				m.activeSlot = cardOne
			case cardThree:
				m.next = 1
				m.activeSlot = cardTwo
			case cardFour:
				m.next = 2
				m.activeSlot = cardThree
			case cardFive:
				m.next = 3
				m.activeSlot = cardFour
			case cardSix:
				m.next = 4
				m.activeSlot = cardFive
			}
		}
	case mainModel:
		switch msg.viewState {
		case model.LobbyView:
			m.viewState = model.LobbyView
			m.activeTab = 1
		case model.PlayView:
			m.gameState = model.PlayState
			m.activeTab = 2
		case model.KittyView:
			m.gameState = model.PlayState
			m.activeTab = 3
		case model.HandView:
			m.gameState = model.PlayState
			m.activeTab = 1
		case model.ScoresView:
			break
		case model.GameOverView:
			break
		default:
			break
		}
	}
	return m, tea.Batch(cmds...)
}

func (m mainModel) View() string {
	switch m.viewState {
	case model.LobbyView:
		return viewStyle.Render(m.lobbyView())
	case model.PlayView:
		return viewStyle.Render(m.gameView())
	case model.KittyView:
		//return viewStyle.Render(m.kittyView())
		break
	case model.HandView:
		//return viewStyle.Render(m.handView())
		break
	default:
		return viewStyle.Render(m.lobbyView())
	}
	return viewStyle.Render(m.lobbyView())
}

func (m mainModel) gameView() string {
	doc := strings.Builder{}
	var renderedTabs []string

	for i, t := range m.tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(m.tabs)-1, i == m.activeTab
		if isActive {
			style = activeTabStyle
		} else {
			style = inactiveTabStyle
		}
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "└"
		} else if isLast && !isActive {
			border.BottomRight = "┴"
		}

		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(t))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	row = lipgloss.JoinHorizontal(lipgloss.Bottom, "────────────────────────────────────────────────────────────────────────────┐")
	doc.WriteString(row)
	doc.WriteString("\n")

	var view string
	switch m.activeTab {
	case 0:
		view = m.handView(m.cardsInPlay)
	case 1:
		view = m.handView(m.hand)
	case 2:
		view = m.handView(m.kitty)
	case 3:
		view = m.handleTestView()
	}

	doc.WriteString(windowStyle.Width(100).Render(view))
	return doc.String()
}

func (m mainModel) lobbyView() string {
	return "Lobby View"
}

func (m mainModel) handView(cards []model.Card) string {
	var s string

	cardViews := make([]string, 0)
	for i := 0; i < len(cards); i++ {
		view := fmt.Sprintf("%s : %v", cards[i].Suit, cards[i].Value)

		if slices.Contains(m.selectedCardIds, cards[i].Id) {
			if i == m.next {
				cardViews = append(cardViews, selectedFocusedStyle.Render(view))
			} else {
				cardViews = append(cardViews, selectedStyle.Render(view))
			}
		} else {
			if i == m.next {
				cardViews = append(cardViews, focusedModelStyle.Render(view))
			} else {
				cardViews = append(cardViews, modelStyle.Render(view))
			}
		}
	}

	s += lipgloss.JoinHorizontal(lipgloss.Top, cardViews...)
	s += helpStyle.Render("\ntab: focus next • q: exit\n")
	return s
}

func (m mainModel) handleTestView() string {
	return ""
}

func (m mainModel) getCardById(id int) model.Card {
	idx := slices.IndexFunc(m.hand, func(c model.Card) bool {
		return c.Id == id
	})

	if idx == -1 {
		return model.Card{}
	}

	return m.hand[idx]
}

func main() {
	p := tea.NewProgram(newModel())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
