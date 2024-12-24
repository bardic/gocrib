package controller

import (
	"context"
	"encoding/json"
	"time"

	"github.com/bardic/gocrib/queries/queries"

	conn "github.com/bardic/gocrib/server/db"
	"github.com/bardic/gocrib/vo"
)

// Deal deals cards to players in a given match
// Return the deck after dealing
func Deal(match *vo.GameMatch) (*vo.GameDeck, error) {
	deck, err := GetDeckByMatchId(match.Deckid)

	if err != nil {
		return nil, err
	}

	cardsPerHand := 6
	if len(match.Players) == 3 {
		cardsPerHand = 5
	}

	for i := 0; i < len(match.Players)*cardsPerHand; i++ {
		var cardId int32
		cardId = deck.Cards[0].Card.ID
		deck.Cards = deck.Cards[1:]
		idx := len(match.Players) - 1 - i%len(match.Players)

		if len(match.Players[idx].Hand) < cardsPerHand {
			match.Players[idx].Hand = append(match.Players[idx].Hand, cardId)
		}
	}

	for _, p := range match.Players {
		_, err := UpdatePlayerById(p)

		if err != nil {
			return nil, err
		}
	}

	return deck, nil
}

// RemoveCardsFromHand takes a hand modifier object and updates the state of the game with that information
func RemoveCardsFromHand(details vo.HandModifier) (*vo.GameMatch, error) {
	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	q.RemoveCardsFromHand(ctx, queries.RemoveCardsFromHandParams{
		ID:   details.PlayerId,
		Hand: details.CardIds,
	})

	b, err := q.GetMatchById(ctx, details.MatchId)

	if err != nil {
		return nil, err
	}

	var match vo.GameMatch
	err = json.Unmarshal(b, &match)
	if err != nil {
		return nil, err
	}

	return &match, nil
}

// UpdatePlay updates the state of the game with the play information
func UpdatePlay(details vo.HandModifier) (*vo.GameMatch, error) {
	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	q.UpdateCardsPlayed(ctx, queries.UpdateCardsPlayedParams{
		Play: details.CardIds,
		ID:   details.PlayerId,
	})

	return RemoveCardsFromHand(details)
}
