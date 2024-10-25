package utils

import (
	"context"
	"encoding/json"

	"model"
	"queries"
	conn "server/db"
)

func UpdateGameState(matchId int32, state queries.Gamestate) error {
	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx := context.Background()

	err := q.UpdateMatchState(ctx, queries.UpdateMatchStateParams{Gamestate: state, ID: matchId})

	if err != nil {
		return err
	}

	return nil
}

func GetMatch(id int) (*model.GameMatch, error) {
	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx := context.Background()

	m, err := q.GetMatchById(ctx, int32(id))

	if err != nil {
		return nil, err
	}

	var match *model.GameMatch
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

	ctx := context.Background()

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

func NewDeck() (queries.Deck, error) {
	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx := context.Background()

	cards, err := q.GetCards(ctx)

	if err != nil {
		return queries.Deck{}, err
	}

	//get card ids
	var cardIds []int32
	for _, card := range cards {
		cardIds = append(cardIds, card.ID)
	}

	deck, err := q.CreateDeck(ctx, cardIds)

	if err != nil {
		return queries.Deck{}, err
	}

	return deck, nil
}

func UpdateMatch(match queries.Match) error {
	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx := context.Background()

	err := q.UpdateMatch(ctx, queries.UpdateMatchParams{
		ID:                 match.ID,
		Playerids:          match.Playerids,
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

func UpdatePlayersInMatch(req model.JoinMatchReq) error {
	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx := context.Background()

	err := q.UpdatePlayersInMatch(ctx, queries.UpdatePlayersInMatchParams{
		ID:          int32(req.MatchId),
		ArrayAppend: int32(req.PlayerId),
	})

	if err != nil {
		return err
	}

	return nil
}
