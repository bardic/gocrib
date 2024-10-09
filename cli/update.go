package main

import (
	"encoding/json"
	"fmt"

	"github.com/bardic/gocrib/cli/services"
	"github.com/bardic/gocrib/cli/state"
	"github.com/bardic/gocrib/cli/utils"
	"github.com/bardic/gocrib/cli/views"
	"github.com/bardic/gocrib/model"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

		if m.gameState == model.CutState {
			var cmd tea.Cmd
			//views.CutInput.Focus()
			views.CutInput, cmd = views.CutInput.Update(msg)
			cmds = append(cmds, cmd)
		}
	case timer.TickMsg:
		if m.ViewStateName != views.Game {
			break
		}

		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		cmds = append(cmds, cmd, services.GetPlayerMatch)
	case model.GameMatch:
		fmt.Println("asdasdasd")
	case int:
		switch m.ViewStateName {
		case views.Login:

		case views.Lobby:
			m.ViewStateName = views.Game
			return m, tea.Batch(cmds...)
		case views.Game:
			if !m.timerStarted {
				m.timerStarted = true
				cmd := m.timer.Init()
				cmds = append(cmds, cmd)
			}
			switch m.gameState {
			case model.WaitingState:
			case model.MatchReady:
			case model.DealState:
			case model.CutState:
			case model.DiscardState:
			case model.PlayState:
			case model.OpponentState:
			case model.KittyState:
			case model.GameWonState:
			case model.GameLostState:

			}
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
			cmds = append(cmds, services.GetPlayerMatch)
			return m, tea.Batch(cmds...)
		case views.Game:
			if !m.timerStarted {
				m.timerStarted = true
				cmd := m.timer.Init()
				cmds = append(cmds, cmd)
			}

			var match model.GameMatch
			err := json.Unmarshal([]byte(msg), &match)

			if err != nil {
				utils.Logger.Info(err.Error())
			}

			var cmd tea.Cmd
			//views.CutInput.Focus()
			views.CutInput, cmd = views.CutInput.Update(msg)
			cmds = append(cmds, cmd)

			if m.gameState == match.GameState {
				return m, tea.Batch(cmds...)
			}

			state.ActiveMatch = match
			m.gameState = match.GameState

			switch m.gameState {
			case model.NewGameState:
			case model.WaitingState:
				fmt.Println("Waiting State")
				deckByte := services.GetDeckById(match.DeckId).([]byte)
				var deck model.GameDeck
				json.Unmarshal(deckByte, &deck)
				state.ActiveDeck = &deck
			case model.MatchReady:
				fmt.Println("Match Ready")
			case model.DealState:
				fmt.Println("Deal state")
			case model.CutState:
				var cmd tea.Cmd
				//views.CutInput.Focus()
				views.CutInput, cmd = views.CutInput.Update(msg)
				cmds = append(cmds, cmd)
				// fmt.Println("Cut state")
			case model.DiscardState:
			case model.PlayState:
			case model.OpponentState:
			case model.KittyState:
			case model.GameWonState:
			case model.GameLostState:

			}
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

func (m *AppModel) setCards(match model.GameMatch) {
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
