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

		cmds = append(cmds, func() tea.Msg {
			return resp
		})
	case model.Account:
		m.accountId = msg.Id
		cmds = append(cmds, func() tea.Msg {
			return model.StateChangeMsg{
				NewState: model.LobbyView,
			}
		})
	case model.MatchDetailsResponse:
		matchDetails := msg
		m.matchId = matchDetails.MatchId

		if matchDetails.GameState == model.NewGameState ||
			matchDetails.GameState == model.JoinGameState {
			m.ViewStateName = model.JoinGameView
			cmds = append(cmds, func() tea.Msg {
				return model.GameStateChangeMsg{
					NewState: model.JoinGameState,
					PlayerId: matchDetails.PlayerId,
					MatchId:  m.matchId,
				}
			})

			break
		}

		m.ViewStateName = model.InGameView
		cmds = append(cmds, func() tea.Msg {
			return model.GameStateChangeMsg{
				NewState: matchDetails.GameState,
				PlayerId: matchDetails.PlayerId,
				MatchId:  m.matchId,
			}
		})
	case model.GameStateChangeMsg:
		switch msg.NewState {
		case model.NewGameState:
			var cmd tea.Cmd
			cmd = m.createMatch(msg)
			gameView := m.currentView.(*views.GameView)
			p, err := utils.GetPlayerId(m.accountId, gameView.GameMatch.Players)
			if err != nil {
				utils.Logger.Sugar().Error(err)
			}
			services.PlayerReady(p.Id)
			cmds = append(cmds, cmd)
		case model.JoinGameState:
			var cmd tea.Cmd
			cmd = m.createMatch(msg)
			services.PlayerReady(msg.PlayerId)
			cmds = append(cmds, cmd)
		case model.WaitingState:
			m.playersReady = true
			m.ViewStateName = model.InGameView
		case model.CutState:
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

		var matchDetails model.GameMatch
		idstr := strconv.Itoa(m.matchId)
		resp := services.GetPlayerMatch(idstr)
		json.Unmarshal(resp.([]byte), &matchDetails)

		gameView := m.currentView.(*views.GameView)
		gameView.GameMatch = &matchDetails
		m.setCards(gameView.GameMatch)

		cmds = append(cmds, cmd, func() tea.Msg {
			return matchDetails
		})
	}

	m.currentView.Update(msg)

	return m, tea.Batch(cmds...)
}

func (m *AppModel) createMatch(msg model.GameStateChangeMsg) tea.Cmd {
	if m.GameInitd {
		return nil
	}

	m.currentView = &views.GameView{
		AccountId: m.accountId,
		MatchId:   msg.MatchId,
	}

	gameView := m.currentView.(*views.GameView)
	gameView.Init()
	m.GameInitd = true
	m.ViewStateName = model.InGameView

	cmd := m.timer.Init()

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
