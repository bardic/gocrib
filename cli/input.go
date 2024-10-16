package main

import (
	"slices"

	"github.com/bardic/gocrib/cli/services"
	"github.com/bardic/gocrib/cli/utils"
	"github.com/bardic/gocrib/cli/views"
	"github.com/bardic/gocrib/model"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *AppModel) parseInput(msg tea.KeyMsg) tea.Msg {
	switch msg.String() {
	case "ctrl+c", "q":
		return nil
	case "enter", "view_update":
		return m.currentView.Enter()
	case "n":
		if m.ViewStateName != model.LobbyView {
			return nil
		}
		m.ViewStateName = model.CreateGameView
		return utils.CreateGame(m.accountId)
	case " ":
		if m.ViewStateName != model.InGameView {
			return nil
		}

		gameView := m.currentView.(*views.GameView)
		cards := m.getVisibleCards(gameView.ActiveTab, gameView.GameMatch.Players[0])

		if len(cards) == 0 {
			return nil
		}

		idx := slices.Index(gameView.HighlightedIds, cards[gameView.HighlighedId])
		if idx > -1 {
			gameView.HighlightedIds = slices.Delete(gameView.HighlightedIds, 0, 1)
		} else {
			gameView.HighlightedIds = append(gameView.HighlightedIds, cards[gameView.HighlighedId])
		}
	case "tab":
		switch m.ViewStateName {
		case model.LobbyView:
			lobbyView := m.currentView.(*views.LobbyView)

			lobbyView.ActiveLandingTab = lobbyView.ActiveLandingTab + 1

			switch lobbyView.ActiveLandingTab {
			case 0:
				lobbyView.LobbyViewState = model.OpenMatches
			case 1:
				lobbyView.LobbyViewState = model.AvailableMatches

			}
		case model.InGameView:
			gameView := m.currentView.(*views.GameView)
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
		switch m.ViewStateName {
		case model.LobbyView:
			lobbyView := m.currentView.(*views.LobbyView)
			lobbyView.ActiveLandingTab = lobbyView.ActiveLandingTab - 1

			switch lobbyView.ActiveLandingTab {
			case 0:
				lobbyView.LobbyViewState = model.OpenMatches
			case 1:
				lobbyView.LobbyViewState = model.AvailableMatches

			}
		case model.InGameView:
			gameView := m.currentView.(*views.GameView)
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

		if m.ViewStateName != model.InGameView {
			return nil
		}

		gameView := m.currentView.(*views.GameView)

		gameView.ActiveSlot++

		cards := m.getVisibleCards(gameView.ActiveTab, gameView.GameMatch.Players[0])

		if gameView.ActiveSlot > len(cards)-1 {
			gameView.ActiveSlot = 0
		}

		gameView.HighlighedId = gameView.ActiveSlot
	case "left":
		if m.ViewStateName != model.InGameView {
			return nil
		}

		gameView := m.currentView.(*views.GameView)
		gameView.ActiveSlot--

		cards := m.getVisibleCards(gameView.ActiveTab, gameView.GameMatch.Players[0])

		if gameView.ActiveSlot < 0 {
			gameView.ActiveSlot = len(cards) - 1
		}

		gameView.HighlighedId = gameView.ActiveSlot
	}

	return nil
}

func (m *AppModel) getVisibleCards(activeTab int, player model.Player) []int {
	var cards []int
	switch activeTab {
	case 0:
		cards = nil
	case 1:
		cards = player.Play
	case 2:
		cards = player.Hand
	case 3:
		cards = player.Kitty
	}

	return cards
}

func (m *AppModel) OnEnterDuringPlay() {
	if m.ViewStateName != model.InGameView {
		return
	}
	gameView := m.currentView.(*views.GameView)

	if gameView.GameState == model.WaitingState {
		gameView.GameState = model.DiscardState
	}

	m.ViewStateName = model.InGameView

	if gameView.GameState == model.DiscardState {
		for _, idx := range gameView.HighlightedIds {
			gameView.Kitty = append(gameView.Kitty, utils.GetCardInHandById(idx, gameView.Hand))
			gameView.Hand = slices.DeleteFunc(gameView.Hand, func(c model.Card) bool {
				return c.Id == idx
			})
		}

		gameView.HighlightedIds = []int{}

		services.PutKitty(model.HandModifier{
			MatchId:  gameView.GameMatch.Id,
			CardIds:  utils.GetIdsFromCards(gameView.Kitty),
			PlayerId: gameView.GameMatch.PlayerIds[0],
		})
	} else {
		for _, idx := range gameView.HighlightedIds {
			gameView.Play = append(gameView.Play, utils.GetCardInHandById(idx, gameView.Hand))
			gameView.Hand = slices.DeleteFunc(gameView.Hand, func(c model.Card) bool {
				return c.Id == idx
			})
		}

		gameView.HighlightedIds = []int{}
		services.PutPlay(model.HandModifier{
			MatchId:  gameView.GameMatch.Id,
			CardIds:  utils.GetIdsFromCards(gameView.Play),
			PlayerId: gameView.GameMatch.PlayerIds[0],
		})
	}
}
