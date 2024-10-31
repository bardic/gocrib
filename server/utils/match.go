package utils

import (
	"context"
	"encoding/json"
	"queries"
	conn "server/db"
	"time"
	"vo"
)

func UpdateGameState(matchId int32, state queries.Gamestate) error {
	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := q.UpdateMatchState(ctx, queries.UpdateMatchStateParams{Gamestate: state, ID: matchId})

	if err != nil {
		return err
	}

	return nil
}

func GetMatch(id int) (*vo.GameMatch, error) {
	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	m, err := q.GetMatchById(ctx, int32(id))

	if err != nil {
		return nil, err
	}

	var match *vo.GameMatch
	err = json.Unmarshal(m, &match)
	if err != nil {
		return nil, err
	}
	return match, nil
}

func GetOpenMatches() ([]queries.Match, error) {
	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	matchesData, err := q.GetOpenMatches(ctx)

	if err != nil {
		return nil, err
	}

	var matches []queries.Match

	for _, matchData := range matchesData {
		var match queries.Match
		err = json.Unmarshal(matchData, &match)
		if err != nil {
			return nil, err
		}

		matches = append(matches, match)
	}

	return matches, nil
}

func GetGameDeck(matchId int32) (*vo.GameDeck, error) {
	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//create a new deck object
	//create match cards based on queries.Cards
	//update deck_matchcard with new match card ids

	deck, err := q.GetDeckForMatchId(ctx, matchId)

	if err != nil {
		return &vo.GameDeck{}, err
	}

	cards, err := q.GetCards(ctx)

	if err != nil {
		return &vo.GameDeck{}, err
	}

	matchcards := []*vo.GameCard{}
	for _, card := range cards {
		matchCard, err := q.CreateMatchCards(ctx, queries.CreateMatchCardsParams{
			Cardid: card.ID,
			State:  queries.CardstateDeck,
		},
		)

		if err != nil {
			return &vo.GameDeck{}, err
		}

		gameCard := &vo.GameCard{
			Matchcard: matchCard,
			Card:      card,
		}

		matchcards = append(matchcards, gameCard)
	}

	gameDeck := &vo.GameDeck{
		Deck:  &deck,
		Cards: matchcards,
	}

	return gameDeck, nil
}

func UpdateMatch(match queries.Match) error {
	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := q.UpdateMatch(ctx, queries.UpdateMatchParams{
		ID:                 match.ID,
		Creationdate:       match.Creationdate,
		Privatematch:       match.Privatematch,
		Elorangemin:        match.Elorangemin,
		Elorangemax:        match.Elorangemax,
		Deckid:             match.Deckid,
		Cutgamecardid:      match.Cutgamecardid,
		Currentplayerturn:  match.Currentplayerturn,
		Turnpasstimestamps: match.Turnpasstimestamps,
		Gamestate:          match.Gamestate,
		Art:                match.Art,
	})

	if err != nil {
		return err
	}

	return nil
}
