package utils

import (
	"context"
	"encoding/json"

	"github.com/bardic/gocrib/model"
	"github.com/bardic/gocrib/queries"
	conn "github.com/bardic/gocrib/server/db"
)

// var matchQuery = `SELECT
// 			json_build_object(
// 				'id', id,
// 				'playerIds', playerIds,
// 				'creationDate', creationDate,
// 				'privateMatch', privateMatch,
// 				'eloRangeMin', eloRangeMin,
// 				'eloRangeMax', eloRangeMax,
// 				'deckid', deckid,
// 				'cutgamecardid', cutgamecardid,
// 				'currentplayerturn', currentplayerturn,
// 				'turnpasstimestamps', turnpasstimestamps,
// 				'art', art,
// 				'gameState', gameState,
// 				'players', (SELECT json_agg(
// 					json_build_object(
// 						'id', p.id,
// 						'accountid', p.accountid,
// 						'play', p.play,
// 						'hand', p.hand,
// 						'kitty', p.kitty,
// 						'score', p.score,
// 						'isready', p.isready,
// 						'art', p.art ))
// 				FROM player as p WHERE p.Id = ANY(m.playerIds)))`

func CardsInPlay(players []queries.Player) []int32 {
	play := []int32{}
	for _, player := range players {
		play = append(play, player.Play...)
	}

	return play
}

func UpdateGameState(matchId int, state queries.Gamestate) error {
	// args := pgx.NamedArgs{
	// 	"id":        matchId,
	// 	"gameState": state,
	// }
	// query := `UPDATE match SET
	// 			gameState=@gameState
	// 		WHERE id=@id`

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
	// db := conn.Pool()
	// defer db.Close()
	// q := queries.New(db)

	// ctx := context.Background()

	// cards, err := q.GetCards(ctx)

	// if err != nil {
	// 	return queries.Deck{}, err
	// }

	// var cardsBytes [][]byte

	// deck, err := q.CreateDeck(ctx)
	// db := conn.Pool()
	// defer db.Close()

	// //rows, err := db.Query(context.Background(), "SELECT * FROM cards") //MEOWCAKES

	// v := []queries.Card{}

	// for rows.Next() {
	// 	var card queries.Card

	// 	err := rows.Scan(&card.Id, &card.Value, &card.Suit, &card.Art)
	// 	if err != nil {
	// 		return model.GameDeck{}, err
	// 	}

	// 	v = append(v, card)
	// }

	// if err != nil {
	// 	return model.GameDeck{}, err
	// }

	// deck := model.GameDeck{
	// 	Cards: []queries.Card{},
	// }

	// for _, c := range v {
	// 	deck.Cards = append(deck.Cards, queries.Card{
	// 		CardId: c.Id,
	// 		Card:   c,
	// 		State:  0,
	// 	})
	// }

	// b, err := json.Marshal(deck.Cards)

	// if err != nil {
	// 	return model.GameDeck{}, err
	// }

	// args := pgx.NamedArgs{"cards": string(b)}

	// //MEOWCAKES 	query := "INSERT INTO deck(cards) VALUES (@cards) RETURNING id"

	// var deckId int
	// err = db.QueryRow(
	// 	context.Background(),
	// 	query,
	// 	args).Scan(&deckId)

	// if err != nil {
	// 	return model.GameDeck{}, err
	// }

	// deck.Id = deckId

	return queries.Deck{}, nil
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

	err := q.UpdatePlayersInMatch(ctx, queries.UpdatePlayersInMatchParams{
		ID:          int32(req.MatchId),
		ArrayAppend: &req.PlayerId,
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

func GetDeckById(id int) (queries.Deck, error) {

	//TODO : DECKS NEED TO BE REIMPLEMENTED
	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx := context.Background()

	d, err := q.GetDeck(ctx, int32(id))

	if err != nil {
		return queries.Deck{}, err
	}

	return d, nil
}
