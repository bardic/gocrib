package main

import (
	"encoding/json"
	"fmt"

	"github.com/bardic/cribbagev2/cli/services"
	"github.com/bardic/cribbagev2/cli/state"
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
		var matchStr string

		err := json.Unmarshal(msg, &matchStr)

		if err != nil {
			fmt.Println(err)
		}

		switch m.ViewStateName {
		case views.Login:
			m.ViewStateName = views.Lobby
			var account model.Account
			err = json.Unmarshal([]byte(matchStr), &account)

			if err != nil {
				fmt.Println(err)
			}

			state.PlayerId = account.Id

			return m, tea.Batch(cmds...)
		case views.Lobby:
			m.ViewStateName = views.Game
			return m, tea.Batch(cmds...)
		}

		var match model.GameMatch
		err = json.Unmarshal([]byte(matchStr), &match)

		if err != nil {
			fmt.Println(err)
		}

		diffs := match.Eq(state.ActiveMatch.Match)

		if diffs == 0 {
			return m, tea.Batch(cmds...)
		}

		state.ActiveMatch = match

		for diff := model.GenericDiff; diff < model.MaxDiff; diff <<= 1 {
			d := diffs & diff
			if d != 0 {
				switch d {
				case model.CutDiff:
					fmt.Println("cutdiff")
					m.ViewStateName = views.Game
				case model.CardsInPlayDiff:
					fmt.Println("cards in play diff")
				case model.GameStateDiff:
					fmt.Println("game state diff")
				case model.GenericDiff:
					fmt.Println("generic diff")
				case model.NewDeckDiff:
					m.ViewStateName = views.Game
					fmt.Println("new deck diff")
					deckByte := services.GetDeckById(match.DeckId).([]byte)
					var deckJson string
					json.Unmarshal(deckByte, &deckJson)

					var deck model.GameDeck
					json.Unmarshal([]byte(deckJson), &deck)

					state.ActiveDeck = deck
					state.ActiveMatch = match
				case model.MaxDiff:
					fmt.Println("max diff")
				case model.TurnDiff:
					fmt.Println("turn diff")
				case model.TurnPassTimestampsDiff:
					fmt.Println("pass timestamp diff")
				}
			}
		}

		m.setCards(match)

		m.gameState = match.GameState
		state.ActiveMatch = match
	}

	return m, tea.Batch(cmds...)
}

func (m *appModel) setCards(match model.GameMatch) {
	m.hand = []model.Card{}
	m.kitty = []model.Card{}
	m.play = []model.Card{}

	for _, cardId := range match.Players[0].Hand {
		card := getCardById(cardId)
		if card != nil {
			m.hand = append(m.hand, *card)
		}
	}

	for _, cardId := range match.Players[0].Kitty {
		card := getCardById(cardId)
		m.kitty = append(m.kitty, *card)
	}

	for _, cardId := range match.Players[0].Play {
		card := getCardById(cardId)
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
