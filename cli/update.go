package main

import (
	"encoding/json"

	"github.com/bardic/cribbagev2/cli/services"
	"github.com/bardic/cribbagev2/cli/state"
	"github.com/bardic/cribbagev2/cli/utils"
	"github.com/bardic/cribbagev2/cli/views"
	"github.com/bardic/cribbagev2/model"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		cmd := m.parseInput(msg)

		if cmd != nil {
			return m, cmd
		}

		if m.ViewStateName == views.Login {
			views.LoginIdField.Focus()
			views.LoginIdField, cmd = views.LoginIdField.Update(msg)
		}

		if m.ViewStateName == views.Lobby {
			views.LobbyTable.Focus()
			views.LobbyTable, cmd = views.LobbyTable.Update(msg)

		}

		cmds = append(cmds, cmd)
	case timer.TickMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)

		if m.ViewStateName == views.Lobby || m.ViewStateName == views.Login {
			cmds = append(cmds, cmd)
		} else {
			cmds = append(cmds, cmd, services.GetPlayerMatch)
		}
	case []byte:
		switch m.ViewStateName {
		case views.Login:
			m.ViewStateName = views.Lobby
			var account model.Account
			err := json.Unmarshal([]byte(msg), &account)

			if err != nil {
				utils.Logger.Info(err.Error())
			}

			state.AccountId = account.Id

			return m, tea.Batch(cmds...)
		case views.Lobby:
			m.ViewStateName = views.Game
			return m, tea.Batch(cmds...)
		}

		var match model.GameMatch
		err := json.Unmarshal([]byte(msg), &match)

		if err != nil {
			utils.Logger.Info(err.Error())
		}

		diffs := match.Eq(state.ActiveMatch)

		if diffs == 0 {
			return m, tea.Batch(cmds...)
		}

		state.ActiveMatch = match
		m.gameState = match.GameState

		for diff := model.GenericDiff; diff < model.MaxDiff; diff <<= 1 {
			d := diffs & diff
			if d != 0 {
				switch d {
				case model.CutDiff:
					utils.Logger.Info("cutdiff")
					m.ViewStateName = views.Game
				case model.CardsInPlayDiff:
					utils.Logger.Info("cards in play diff")
				case model.GameStateDiff:
					utils.Logger.Info("game state diff")
				case model.GenericDiff:
					// utils.Logger.Info("generic diff")
					m.ViewStateName = views.Game
					// utils.Logger.Info("new deck diff")
					// deckByte := services.GetDeckById(match.DeckId).([]byte)
					// var deck model.GameDeck
					// json.Unmarshal(deckByte, &deck)
					// state.ActiveDeck = &deck
				case model.NewDeckDiff:
					m.ViewStateName = views.Game
					utils.Logger.Info("new deck diff")
					deckByte := services.GetDeckById(match.DeckId).([]byte)
					var deck model.GameDeck
					json.Unmarshal(deckByte, &deck)
					state.ActiveDeck = &deck
				case model.MaxDiff:
					utils.Logger.Info("max diff")
				case model.PlayersDiff:
					utils.Logger.Info("players diff")
				case model.TurnDiff:
					utils.Logger.Info("turn diff")
				case model.TurnPassTimestampsDiff:
					utils.Logger.Info("pass timestamp diff")
				}
			}
		}

		if state.ActiveDeck != nil {
			utils.Logger.Info("set cards")
			m.setCards(match)
		}
	}

	return m, tea.Batch(cmds...)
}

func getPlayerForId(id int, match model.GameMatch) *model.Player {
	for _, player := range match.Players {
		if player.AccountId == id {
			return &player
		}
	}

	return nil
}

func (m *appModel) setCards(match model.GameMatch) {
	m.hand = []model.Card{}
	m.kitty = []model.Card{}
	m.play = []model.Card{}

	p := getPlayerForId(state.AccountId, match)

	for _, cardId := range p.Hand {
		card := utils.GetCardById(cardId)
		if card != nil {
			m.hand = append(m.hand, *card)
		}
	}

	for _, cardId := range p.Kitty {
		card := utils.GetCardById(cardId)
		m.kitty = append(m.kitty, *card)
	}

	for _, cardId := range p.Play {
		card := utils.GetCardById(cardId)
		m.play = append(m.play, *card)
	}
}

func createGame() tea.Msg {
	newMatch := services.PostPlayerMatch().([]byte)

	var match model.GameMatch
	json.Unmarshal(newMatch, &match)
	state.ActiveMatchId = match.Id
	state.ActiveMatch = match
	return match
}
