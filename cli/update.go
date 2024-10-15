package main

import (
	"encoding/json"

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

	//Parse msg
	switch msg := msg.(type) {
	case tea.KeyMsg: //User input
		if cmd := m.parseInput(msg); cmd != nil {
			cmds = append(cmds, cmd)
		}
		m.currentView.(*views.LoginView).LoginIdField.Focus()
	case timer.TickMsg: // Polling update
		if state.ViewStateName != model.GameView {
			break
		}

		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)

		cmds = append(cmds, cmd, services.GetPlayerMatchState)
	case []byte: // Response from HTTP service
		if cmd, err := m.updateView(msg); err != nil {
			utils.Logger.Sugar().Error(err)
		} else {
			cmds = append(cmds, cmd)
		}
	}

	m.currentView.Update(msg)

	return m, tea.Batch(cmds...)
}

// Responsible for updating View State and
func (m *AppModel) updateView(msg []byte) (tea.Cmd, error) {
	var cmds []tea.Cmd
	switch state.ViewStateName {
	case model.LoginView:
		var account model.Account
		if err := json.Unmarshal(msg, &account); err != nil {
			utils.Logger.Info(err.Error())
		}

		state.AccountId = account.Id
	case model.LobbyView:
		services.GetOpenMatches()
		m.currentView = views.LobbyView{}

		//cmds = append(cmds, services.GetPlayerMatchState)
	case model.GameView:
		if !m.timerStarted {
			m.timerStarted = true
			cmd := m.timer.Init()
			cmds = append(cmds, cmd)
		}

		m.currentView = views.GameView{
			GameState: m.gameState,
			Cards:     []model.Card{},
		}

		//If the game state has changed, we will receive a model.GameMatch
		if state.MatchDetailsResponse != nil {
			var match *model.GameMatch
			if err := json.Unmarshal([]byte(msg), &match); err != nil {
				return nil, err
			}

			//Update state with new match data
			state.ActiveMatch = match

			//If active deck doesn't exist, get it
			if state.ActiveDeck == nil {
				deckByte := services.GetDeckById(match.DeckId).([]byte)
				var deck model.GameDeck
				json.Unmarshal(deckByte, &deck)
				state.ActiveDeck = &deck
			}

			//Update model with cards from latest match
			m.setCards(match)

			state.MatchDetailsResponse = nil
		}

		//Check match for change in statew
		var resp *model.MatchDetailsResponse
		if err := json.Unmarshal([]byte(msg), &resp); err != nil {
			return nil, err
		}

		//No state change
		if m.gameState == resp.GameState {
			break
		}

		//Update state and response
		m.gameState = resp.GameState
		state.MatchDetailsResponse = resp

		//Get updated playermatch
		cmds = append(cmds, services.GetPlayerMatch)
	}

	cmd := m.currentView.Update(msg)

	cmds = append(cmds, cmd)

	return tea.Batch(cmds...), nil
}

func (m *AppModel) setCards(match *model.GameMatch) {
	m.hand = []model.Card{}
	m.kitty = []model.Card{}
	m.play = []model.Card{}

	p := utils.GetPlayerForId(state.AccountId, match)

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
