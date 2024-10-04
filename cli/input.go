package main

import (
	"slices"
	"strconv"

	"github.com/bardic/cribbagev2/cli/services"
	"github.com/bardic/cribbagev2/cli/state"
	"github.com/bardic/cribbagev2/cli/views"
	"github.com/bardic/cribbagev2/model"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *appModel) parseInput(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "ctrl+c", "q":
		return tea.Quit
	case "enter", "view_update":
		m.currentView.Enter()

		switch m.ViewStateName {
		case views.Login:
			idStr := views.LoginIdField.Value()
			id, err := strconv.Atoi(idStr)
			if err != nil {
				return tea.Quit
			}

			state.PlayerId = id
			return services.Login
		case views.Lobby:
			idStr := views.LobbyTable.SelectedRow()[0]
			id, err := strconv.Atoi(idStr)
			if err != nil {
				return tea.Quit
			}

			state.ActiveMatchId = id
			return services.JoinMatch
		case views.Game:
			return m.OnEnterDuringPlay()
		case views.GameOver:
		}
	case "n":
		createGame()
		m.ViewStateName = views.Game
	case " ":
		cards := m.hand
		idx := slices.Index(m.HighlightedIds, cards[m.HighlighedId].Id)
		if idx > -1 {
			m.HighlightedIds = slices.Delete(m.HighlightedIds, idx, 1)
		} else {
			m.HighlightedIds = append(m.HighlightedIds, cards[m.HighlighedId].Id)
		}
	case "tab":
		switch m.ViewStateName {
		case views.Lobby:
			m.ActiveLandingTab = m.ActiveLandingTab + 1
			switch m.ActiveLandingTab {
			case 0:
				m.LobbyViewState = model.OpenMatches
			case 1:
				m.LobbyViewState = model.AvailableMatches

			}
		case views.Game:
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
		switch m.ViewStateName {
		case views.Lobby:
			m.ActiveLandingTab = m.ActiveLandingTab - 1
			switch m.ActiveLandingTab {
			case 0:
				m.LobbyViewState = model.OpenMatches
			case 1:
				m.LobbyViewState = model.AvailableMatches

			}
		case views.Game:
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

	return nil
}

func (m *appModel) OnEnterDuringPlay() tea.Cmd {
	if m.gameState == model.WaitingState {
		m.gameState = model.DiscardState
	}

	m.ViewStateName = views.Game
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
		return services.PutKitty
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
		return services.PutPlay
	}
}
