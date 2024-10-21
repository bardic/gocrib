package utils

import (
	"context"
	"encoding/json"

	"github.com/bardic/gocrib/model"
	"github.com/bardic/gocrib/queries"
	conn "github.com/bardic/gocrib/server/db"
)

func CardsInPlay(players []queries.Player) []int32 {
	play := []int32{}
	for _, player := range players {
		play = append(play, player.Play...)
	}

	return play
}

func UpdateGameState(matchId int, state queries.Gamestate) error {
	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx := context.Background()

	err := q.UpdateMatchState(ctx, queries.UpdateMatchStateParams{Gamestate: state, ID: int32(matchId)})

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

func UpdateCut(matchId int, cutCardId int) error {
	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx := context.Background()

	err := q.UpdateMatchCut(ctx, queries.UpdateMatchCutParams{
		ID:            int32(matchId),
		Cutgamecardid: int32(cutCardId),
	})

	if err != nil {
		return err
	}

	return nil
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

func IsMatchReadyToStart(m *queries.Match) (bool, error) {
	if len(m.Playerids) == 2 {
		return true, nil
	}

	return false, nil
}

func UpdateMatchState(matchId int, state queries.Gamestate) error {
	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx := context.Background()

	q.UpdateMatch(ctx, queries.UpdateMatchParams{
		ID:        int32(matchId),
		Gamestate: state,
	})

	return nil
}

func UpdateMatchCut(cardId, matchId int) error {
	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx := context.Background()

	q.UpdateMatchCut(ctx, queries.UpdateMatchCutParams{
		ID:            int32(matchId),
		Cutgamecardid: int32(cardId),
	})

	return nil
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

func UpdatePlayersInMatch(req model.JoinMatchReq) (*model.GameMatch, error) {
	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx := context.Background()

	deck, err := NewDeck()

	if err != nil {
		return nil, err
	}

	err = q.StartMatch(ctx, queries.StartMatchParams{
		ID:          int32(req.MatchId),
		Deckid:      deck.ID,
		ArrayAppend: req.PlayerId,
	})

	if err != nil {
		return nil, err
	}

	m, err := q.GetMatchById(ctx, int32(req.MatchId))

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

func GetDeckById(id int32) (queries.Deck, error) {
	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx := context.Background()

	d, err := q.GetDeck(ctx, id)

	if err != nil {
		return queries.Deck{}, err
	}

	return d, nil
}
