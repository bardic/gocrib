package main

import (
	"slices"

	"github.com/bardic/gocrib/cli/services"
	"github.com/bardic/gocrib/cli/state"
	"github.com/bardic/gocrib/cli/utils"
	"github.com/bardic/gocrib/cli/views"
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
		if state.ViewStateName != model.GameView {
			return nil
		}

		gameView := m.currentView.(views.GameView)
		cards := m.hand
		idx := slices.Index(gameView.HighlightedIds, cards[gameView.HighlighedId].Id)
		if idx > -1 {
			gameView.HighlightedIds = slices.Delete(gameView.HighlightedIds, idx, 1)
		} else {
			gameView.HighlightedIds = append(gameView.HighlightedIds, cards[gameView.HighlighedId].Id)
		}
	case "tab":
		switch state.ViewStateName {
		case model.LobbyView:
			lobbyView := m.currentView.(views.LobbyView)

			lobbyView.ActiveLandingTab = lobbyView.ActiveLandingTab + 1

			switch lobbyView.ActiveLandingTab {
			case 0:
				lobbyView.LobbyViewState = model.OpenMatches
			case 1:
				lobbyView.LobbyViewState = model.AvailableMatches

			}
		case model.GameView:
			gameView := m.currentView.(views.GameView)
			gameView.ActiveTab = gameView.ActiveTab + 1

			switch gameView.ActiveTab {
			case 0:
				gameView.GameViewState = model.BoardView
			case 1:
				gameView.GameViewState = model.PlayView
			case 2:
				gameView.GameViewState = model.HandView
			case 3:
				gameView.GameViewState = model.KittyView
			}
		}
	case "shift+tab":
		switch state.ViewStateName {
		case model.LobbyView:
			lobbyView := m.currentView.(views.LobbyView)
			lobbyView.ActiveLandingTab = lobbyView.ActiveLandingTab - 1

			switch lobbyView.ActiveLandingTab {
			case 0:
				lobbyView.LobbyViewState = model.OpenMatches
			case 1:
				lobbyView.LobbyViewState = model.AvailableMatches

			}
		case model.GameView:
			gameView := m.currentView.(views.GameView)
			gameView.ActiveTab = gameView.ActiveTab - 1

			switch gameView.ActiveTab {
			case 0:
				gameView.GameViewState = model.BoardView
			case 1:
				gameView.GameViewState = model.PlayView
			case 2:
				gameView.GameViewState = model.HandView
			case 3:
				gameView.GameViewState = model.KittyView
			}
		}
	case "right":

		if state.ViewStateName != model.GameView {
			return nil
		}

		gameView := m.currentView.(views.GameView)

		gameView.ActiveSlot++

		if gameView.ActiveSlot > len(state.ActiveMatch.Players[0].Play) {
			gameView.ActiveSlot = 0
		}

		gameView.HighlighedId = gameView.ActiveSlot
	case "left":
		if state.ViewStateName != model.GameView {
			return nil
		}

		gameView := m.currentView.(views.GameView)
		gameView.ActiveSlot--

		if gameView.ActiveSlot < 0 {
			gameView.ActiveSlot = len(state.ActiveMatch.Players[0].Play) - 1
		}

		gameView.HighlighedId = gameView.ActiveSlot
	}

	return nil
}

func (m *AppModel) OnEnterDuringPlay() tea.Cmd {
	if m.gameState == model.WaitingState {
		m.gameState = model.DiscardState
	}

	state.ViewStateName = model.GameView

	gameView := m.currentView.(views.GameView)

	if m.gameState == model.DiscardState {
		for _, idx := range gameView.HighlightedIds {
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

		gameView.HighlightedIds = []int{}
		return services.PutKitty
	} else {
		for _, idx := range gameView.HighlightedIds {
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

		gameView.HighlightedIds = []int{}
		return services.PutPlay
	}
}
