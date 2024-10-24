package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/bardic/gocrib/cli/services"
	"github.com/bardic/gocrib/cli/utils"
	"github.com/bardic/gocrib/cli/views/game"
	"github.com/bardic/gocrib/cli/views/lobby"
	"github.com/bardic/gocrib/model"
	"github.com/bardic/gocrib/queries"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	//Parse msg
	switch msg := msg.(type) {
	case tea.KeyMsg: //User input
		gameView := m.currentView
		resp := gameView.ParseInput(msg)

		if resp == nil {
			break
		}

		cmds = append(cmds, func() tea.Msg {
			return resp
		})
	case queries.Account:
		m.account = &msg
		cmds = append(cmds, func() tea.Msg {
			return model.StateChangeMsg{
				NewState: model.LobbyView,
			}
		})
	case model.MatchDetailsResponse:
		matchDetails := msg
		m.matchId = matchDetails.MatchId

		if matchDetails.GameState == queries.GamestateNewGameState ||
			matchDetails.GameState == queries.GamestateJoinGameState {
			m.ViewStateName = model.JoinGameView
			cmds = append(cmds, func() tea.Msg {
				return model.GameStateChangeMsg{
					NewState: queries.GamestateJoinGameState,
					PlayerId: int32(matchDetails.PlayerId),
					MatchId:  int32(m.matchId),
				}
			})

			break
		}

		m.ViewStateName = model.InGameView
		cmds = append(cmds, func() tea.Msg {
			return model.GameStateChangeMsg{
				NewState: matchDetails.GameState,
				PlayerId: int32(matchDetails.PlayerId),
				MatchId:  int32(m.matchId),
			}
		})
	case model.GameStateChangeMsg:
		switch msg.NewState {
		case queries.GamestateJoinGameState:
			var cmd tea.Cmd
			cmd = m.createMatchView(msg)
			services.PlayerReady(msg.PlayerId)
			cmds = append(cmds, cmd)
		case queries.GamestateWaitingState:
			m.playersReady = true
			m.ViewStateName = model.InGameView
		case queries.GamestateDiscardState:
			m.playersReady = true
			m.ViewStateName = model.InGameView
			fmt.Println("Discard State")
		case queries.GamestateCutState:
			m.playersReady = true
			m.ViewStateName = model.InGameView
		}

	case model.StateChangeMsg:
		switch msg.NewState {
		case model.LobbyView:
			m.currentView = &lobby.LobbyView{
				AccountId: m.account.ID,
			}
			m.ViewStateName = model.LobbyView
			services.GetOpenMatches()
		}

	case timer.TickMsg: // Polling update
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)

		if m.ViewStateName != model.InGameView {
			cmds = append(cmds, cmd)
			break
		}

		var match queries.Match
		idstr := strconv.Itoa(m.matchId)
		resp := services.GetPlayerMatch(idstr)
		err := json.Unmarshal(resp.([]byte), &match)
		if err != nil {
			utils.Logger.Sugar().Error(err)
		}

		gameView := m.currentView.(*game.GameView)

		if gameView.GameMatch != nil {
			if match.Gamestate == gameView.GameMatch.Gamestate {
				cmds = append(cmds, cmd)
				m.setCards(gameView.GameMatch)
				break
			}
		}

		players := m.getPlayers(match)

		gameMatch := &model.GameMatch{
			Match:   match,
			Players: players,
		}

		p := utils.GetPlayerForAccountId(m.account.ID, gameMatch)
		gameView.LocalPlayer = p

		deckByte := services.GetDeckById(int(match.Deckid)).([]byte)
		var deck queries.Deck
		err = json.Unmarshal(deckByte, &deck)

		if err != nil {
			utils.Logger.Sugar().Error(err)
		}

		var gameCards []queries.GetGameCardsForMatchRow
		cardsMsg := services.GetGampleCardsForMatch(int(match.ID))
		err = json.Unmarshal(cardsMsg.([]byte), &gameCards)

		if err != nil {
			utils.Logger.Sugar().Error(err)
		}

		gameDeck := &model.GameDeck{
			Deck:  deck,
			Cards: gameCards,
		}

		gameView.Deck = gameDeck

		if err != nil {
			utils.Logger.Sugar().Error(err)
		}

		gameView.GameMatch = gameMatch
		if gameMatch.Gamestate == queries.GamestateDiscardState {
			m.setCards(gameView.GameMatch)
		}

		cmds = append(cmds, cmd, func() tea.Msg {
			return gameMatch
		})
	}

	m.currentView.Update(msg)

	return m, tea.Batch(cmds...)
}

func (m *AppModel) getPlayers(match queries.Match) []queries.Player {
	players := []queries.Player{}

	for _, playerId := range match.Playerids {
		//get player for each playerid from services
		resp := services.GetPlayer(playerId)
		var player queries.Player
		err := json.Unmarshal(resp.([]byte), &player)
		if err != nil {
			utils.Logger.Sugar().Error(err)
		}
		players = append(players, player)
	}

	return players
}

func (m *AppModel) createMatchView(msg model.GameStateChangeMsg) tea.Cmd {
	if m.GameInitd {
		return nil
	}

	m.currentView = &game.GameView{
		Account: m.account,
		MatchId: msg.MatchId,
	}

	gameView := m.currentView.(*game.GameView)
	gameView.Init()
	m.GameInitd = true
	m.ViewStateName = model.InGameView

	cmd := m.timer.Init()

	return cmd
}

func (m *AppModel) setCards(match *model.GameMatch) {
	gameView := m.currentView.(*game.GameView)

	p := utils.GetPlayerForAccountId(m.account.ID, match)

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
