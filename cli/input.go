package main

import (
	"slices"

	"github.com/bardic/gocrib/cli/services"
	"github.com/bardic/gocrib/cli/utils"
	"github.com/bardic/gocrib/cli/views"
	"github.com/bardic/gocrib/model"
	"github.com/bardic/gocrib/queries"
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
		cards := utils.GetVisibleCards(gameView.ActiveTab, gameView.GameMatch.Players[0])

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

		cards := utils.GetVisibleCards(gameView.ActiveTab, gameView.GameMatch.Players[0])

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

		cards := utils.GetVisibleCards(gameView.ActiveTab, gameView.GameMatch.Players[0])

		if gameView.ActiveSlot < 0 {
			gameView.ActiveSlot = len(cards) - 1
		}

		gameView.HighlighedId = gameView.ActiveSlot
	}

	return nil
}

func (m *AppModel) OnEnterDuringPlay() {
	if m.ViewStateName != model.InGameView {
		return
	}
	gameView := m.currentView.(*views.GameView)
	m.ViewStateName = model.InGameView

	if gameView.GameMatch.Gamestate == queries.GamestateCutState {
		for _, idx := range gameView.HighlightedIds {
			gameView.Kitty = append(gameView.Kitty, utils.GetCardInHandById(idx, gameView.Hand))
			gameView.Hand = slices.DeleteFunc(gameView.Hand, func(c queries.Card) bool {
				return c.ID == idx
			})
		}

		gameView.HighlightedIds = []int32{}

		services.PutKitty(model.HandModifier{
			MatchId:  gameView.GameMatch.ID,
			CardIds:  utils.GetIdsFromCards(gameView.Kitty),
			PlayerId: gameView.GameMatch.Playerids[0],
		})
	} else {
		for _, idx := range gameView.HighlightedIds {
			gameView.Play = append(gameView.Play, utils.GetCardInHandById(idx, gameView.Hand))
			gameView.Hand = slices.DeleteFunc(gameView.Hand, func(c queries.Card) bool {
				return c.ID == idx
			})
		}

		gameView.HighlightedIds = []int32{}
		services.PutPlay(model.HandModifier{
			MatchId:  gameView.GameMatch.ID,
			CardIds:  utils.GetIdsFromCards(gameView.Play),
			PlayerId: gameView.GameMatch.Playerids[0],
		})
	}
}
