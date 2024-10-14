package main

import (
	"slices"

	"github.com/bardic/gocrib/cli/services"
	"github.com/bardic/gocrib/cli/state"
	"github.com/bardic/gocrib/cli/utils"
	"github.com/bardic/gocrib/model"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *AppModel) parseInput(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "ctrl+c", "q":
		return tea.Quit
	case "enter", "view_update":
		return m.currentView.Enter
	case "n":
		utils.CreateGame()
		state.ViewStateName = model.GameView
	case " ":
		cards := m.hand
		idx := slices.Index(m.HighlightedIds, cards[m.HighlighedId].Id)
		if idx > -1 {
			m.HighlightedIds = slices.Delete(m.HighlightedIds, idx, 1)
		} else {
			m.HighlightedIds = append(m.HighlightedIds, cards[m.HighlighedId].Id)
		}
	case "tab":
		switch state.ViewStateName {
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
		switch state.ViewStateName {
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
		if state.ViewStateName != model.GameView {
			return nil
		}

		m.ActiveSlot++

		if m.ActiveSlot > len(state.ActiveMatch.Players[0].Play) {
			m.ActiveSlot = 0
		}

		m.HighlighedId = m.ActiveSlot
	case "left":
		m.ActiveSlot--

		if m.ActiveSlot < 0 {
			m.ActiveSlot = len(state.ActiveMatch.Players[0].Play) - 1
		}

		m.HighlighedId = m.ActiveSlot
	}

	return nil
}

func (m *AppModel) OnEnterDuringPlay() tea.Cmd {
	if m.gameState == model.WaitingState {
		m.gameState = model.DiscardState
	}

	state.ViewStateName = model.GameView
	if m.gameState == model.DiscardState {
		for _, idx := range m.HighlightedIds {
			m.kitty = append(m.kitty, utils.GetCardInHandById(idx, m.hand))
			m.hand = slices.DeleteFunc(m.hand, func(c model.Card) bool {
				return c.Id == idx
			})
		}

		state.CurrentHandModifier = model.HandModifier{
			MatchId:  state.ActiveMatchId,
			CardIds:  utils.GetIdsFromCards(m.kitty),
			PlayerId: state.ActiveMatch.PlayerIds[0],
		}

		m.HighlightedIds = []int{}
		return services.PutKitty
	} else {
		for _, idx := range m.HighlightedIds {
			m.play = append(m.play, utils.GetCardInHandById(idx, m.hand))
			m.hand = slices.DeleteFunc(m.hand, func(c model.Card) bool {
				return c.Id == idx
			})
		}

		state.CurrentHandModifier = model.HandModifier{
			MatchId:  state.ActiveMatchId,
			CardIds:  utils.GetIdsFromCards(m.play),
			PlayerId: state.ActiveMatch.PlayerIds[0],
		}

		m.HighlightedIds = []int{}
		return services.PutPlay
	}
}
