package controller

import (
	"context"
	"errors"
	"queries"
	"time"

	conn "github.com/bardic/gocrib/server/db"
	"github.com/bardic/gocrib/vo"
)

// GetDeckByMatchId returns the deck for a match
func GetDeckByMatchId(matchId int32) (*vo.GameDeck, error) {
	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	d, err := q.GetMatchCards(ctx, matchId)

	if err != nil {
		return &vo.GameDeck{}, err
	}

	if len(d) == 0 {
		return &vo.GameDeck{}, errors.New("no deck found")
	}

	gameDeck := &vo.GameDeck{
		Cards: []*vo.GameCard{},
		Deck:  &queries.Deck{},
	}

	for _, matchCardsRow := range d {
		gameDeck.Cards = append(gameDeck.Cards, &vo.GameCard{
			Matchcard: matchCardsRow.Matchcard,
			Card:      matchCardsRow.Card,
		})
	}

	gameDeck.Deck = &d[0].Deck

	return gameDeck, nil
}
