package main

import (
	"encoding/json"
	"strconv"

	"github.com/bardic/gocrib/cli/services"
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
		resp := m.parseInput(msg)
		if resp == nil {
			break
		}

		switch resp.(type) {
		case model.Account:
			m.accountId = resp.(model.Account).Id
			cmds = append(cmds, func() tea.Msg {
				return model.StateChangeMsg{
					NewState: model.LobbyView,
				}
			})
		case model.MatchDetailsResponse:
			matchDetails := resp.(model.MatchDetailsResponse)
			m.matchId = matchDetails.MatchId

			if matchDetails.GameState == model.NewGameState {
				m.ViewStateName = model.CreateGameView
				cmds = append(cmds, func() tea.Msg {
					return model.GameStateChangeMsg{
						NewState: model.JoinGameState,
						MatchId:  m.matchId,
					}
				})
			} else if matchDetails.GameState == model.JoinGameState {
				m.ViewStateName = model.JoinGameView
				cmds = append(cmds, func() tea.Msg {
					return model.GameStateChangeMsg{
						NewState: model.WaitingState,
						MatchId:  m.matchId,
					}
				})
			} else {
				m.ViewStateName = model.InGameView
				cmds = append(cmds, func() tea.Msg {
					return model.GameStateChangeMsg{
						NewState: matchDetails.GameState,
						MatchId:  m.matchId,
					}
				})
			}
		}
	case model.GameStateChangeMsg:
		switch m.ViewStateName {
		case model.CreateGameView:
			var cmd tea.Cmd
			cmd = m.createMatch(msg, model.NewGameState)
			m.currentView.(*views.GameView).GameState = msg.NewState
			cmds = append(cmds, cmd)
		case model.JoinGameView:
			var cmd tea.Cmd
			cmd = m.createMatch(msg, model.JoinGameState)
			m.currentView.(*views.GameView).GameState = msg.NewState
			cmds = append(cmds, cmd)
		case model.PlayersReadyView:
			m.playersReady = true
			m.ViewStateName = model.InGameView
		}

	case model.StateChangeMsg:
		switch msg.NewState {
		case model.LobbyView:
			m.currentView = &views.LobbyView{
				AccountId: m.accountId,
			}
			m.ViewStateName = model.LobbyView
			services.GetOpenMatches()
		}
	case timer.TickMsg: // Polling update
		if m.ViewStateName != model.InGameView {
			break
		}

		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)

		var matchDetails *model.MatchDetailsResponse
		id := m.currentView.(*views.GameView).GameMatch.Id
		idstr := strconv.Itoa(id)
		resp := services.GetPlayerMatchState(idstr)
		json.Unmarshal(resp.([]byte), &matchDetails)

		cmds = append(cmds, cmd, func() tea.Msg {
			return matchDetails
		})
	}

	m.currentView.Update(msg)

	return m, tea.Batch(cmds...)
}

func (m *AppModel) createMatch(msg model.GameStateChangeMsg, state model.GameState) tea.Cmd {
	if m.GameInitd {
		return nil
	}

	m.currentView = &views.GameView{
		GameState: state,
		AccountId: m.accountId,
		MatchId:   msg.MatchId,
	}

	gameView := m.currentView.(*views.GameView)
	gameView.Init()
	gameView.UpdateState(model.WaitingState)
	m.GameInitd = true

	cmd := m.timer.Init()

	m.ViewStateName = model.InGameView

	return cmd
}

func (m *AppModel) setCards(match *model.GameMatch) {
	gameView := m.currentView.(*views.GameView)
	gameView.Hand = []model.Card{}
	gameView.Kitty = []model.Card{}
	gameView.Play = []model.Card{}

	p := utils.GetPlayerForAccountId(m.accountId, match)

	for _, cardId := range p.Hand {
		card := utils.GetCardById(cardId, gameView.Deck)
		if card != nil {
			gameView.Hand = append(gameView.Hand, *card)
		}
	}

	for _, cardId := range p.Kitty {
		card := utils.GetCardById(cardId, gameView.Deck)
		gameView.Kitty = append(gameView.Kitty, *card)
	}

	for _, cardId := range p.Play {
		card := utils.GetCardById(cardId, gameView.Deck)
		gameView.Play = append(gameView.Play, *card)
	}
}
