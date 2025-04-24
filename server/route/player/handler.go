package player

import (
	"github.com/bardic/gocrib/server/route/deck"
	"github.com/bardic/gocrib/server/store"
)

type PlayerHandler struct {
	// gameDeck  *vo.GameDeck
	// player    *queries.Player
	// gameHands map[queries.Cardstate][]vo.GameCard

	deckHandler *deck.HandlerPlayerDeck
	playerStore store.PlayerStore
	deckStore   store.DeckStore
	cardStore   store.CardStore
	matchStore  store.MatchStore
}

func NewPlayerHandler() *PlayerHandler {
	return &PlayerHandler{}
}
