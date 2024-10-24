// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package queries

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createDeck = `-- name: CreateDeck :one
INSERT INTO deck(cards) VALUES ($1) RETURNING id, cards
`

func (q *Queries) CreateDeck(ctx context.Context, cards []int32) (Deck, error) {
	row := q.db.QueryRow(ctx, createDeck, cards)
	var i Deck
	err := row.Scan(&i.ID, &i.Cards)
	return i, err
}

const createMatch = `-- name: CreateMatch :one
INSERT INTO match(
				playerIds,
				privateMatch,
				eloRangeMin,
				eloRangeMax,
				deckId,
				cutGameCardId,
				currentplayerturn,
				turnPassTimestamps,
				gameState,
				art)
			VALUES (
				$1,
				$2,
				$3,
				$4,
				$5,
				$6,
				$7,
				$8,
				$9,
				$10)
			RETURNING id, playerids, creationdate, privatematch, elorangemin, elorangemax, deckid, cutgamecardid, currentplayerturn, turnpasstimestamps, gamestate, art
`

type CreateMatchParams struct {
	Playerids          []int32
	Privatematch       bool
	Elorangemin        int32
	Elorangemax        int32
	Deckid             int32
	Cutgamecardid      int32
	Currentplayerturn  int32
	Turnpasstimestamps []pgtype.Timestamptz
	Gamestate          Gamestate
	Art                string
}

func (q *Queries) CreateMatch(ctx context.Context, arg CreateMatchParams) (Match, error) {
	row := q.db.QueryRow(ctx, createMatch,
		arg.Playerids,
		arg.Privatematch,
		arg.Elorangemin,
		arg.Elorangemax,
		arg.Deckid,
		arg.Cutgamecardid,
		arg.Currentplayerturn,
		arg.Turnpasstimestamps,
		arg.Gamestate,
		arg.Art,
	)
	var i Match
	err := row.Scan(
		&i.ID,
		&i.Playerids,
		&i.Creationdate,
		&i.Privatematch,
		&i.Elorangemin,
		&i.Elorangemax,
		&i.Deckid,
		&i.Cutgamecardid,
		&i.Currentplayerturn,
		&i.Turnpasstimestamps,
		&i.Gamestate,
		&i.Art,
	)
	return i, err
}

const createMatchCards = `-- name: CreateMatchCards :one
INSERT INTO matchcards (match_id, cardid, origowner, currowner, state) VALUES ($1, $2, $3, $4, $5) RETURNING id, match_id, cardid, origowner, currowner, state
`

type CreateMatchCardsParams struct {
	MatchID   int32
	Cardid    int32
	Origowner pgtype.Int4
	Currowner pgtype.Int4
	State     Cardstate
}

func (q *Queries) CreateMatchCards(ctx context.Context, arg CreateMatchCardsParams) (Matchcard, error) {
	row := q.db.QueryRow(ctx, createMatchCards,
		arg.MatchID,
		arg.Cardid,
		arg.Origowner,
		arg.Currowner,
		arg.State,
	)
	var i Matchcard
	err := row.Scan(
		&i.ID,
		&i.MatchID,
		&i.Cardid,
		&i.Origowner,
		&i.Currowner,
		&i.State,
	)
	return i, err
}

const createPlayer = `-- name: CreatePlayer :one
INSERT INTO player (
			accountid,
			hand,
			play,
			kitty,
			score,
			isready,
			art
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7
		)
		RETURNING id, accountid, play, hand, kitty, score, isready, art
`

type CreatePlayerParams struct {
	Accountid int32
	Hand      []int32
	Play      []int32
	Kitty     []int32
	Score     int32
	Isready   bool
	Art       string
}

func (q *Queries) CreatePlayer(ctx context.Context, arg CreatePlayerParams) (Player, error) {
	row := q.db.QueryRow(ctx, createPlayer,
		arg.Accountid,
		arg.Hand,
		arg.Play,
		arg.Kitty,
		arg.Score,
		arg.Isready,
		arg.Art,
	)
	var i Player
	err := row.Scan(
		&i.ID,
		&i.Accountid,
		&i.Play,
		&i.Hand,
		&i.Kitty,
		&i.Score,
		&i.Isready,
		&i.Art,
	)
	return i, err
}

const getAccount = `-- name: GetAccount :one
SELECT accounts.id, accounts.name FROM accounts WHERE id = $1 LIMIT 1
`

func (q *Queries) GetAccount(ctx context.Context, id int32) (Account, error) {
	row := q.db.QueryRow(ctx, getAccount, id)
	var i Account
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const getCards = `-- name: GetCards :many
SELECT cards.id, cards.value, cards.suit, cards.art FROM cards
`

func (q *Queries) GetCards(ctx context.Context) ([]Card, error) {
	rows, err := q.db.Query(ctx, getCards)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Card
	for rows.Next() {
		var i Card
		if err := rows.Scan(
			&i.ID,
			&i.Value,
			&i.Suit,
			&i.Art,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getCurrentPlayerTurn = `-- name: GetCurrentPlayerTurn :one
SELECT currentplayerturn FROM match WHERE id = $1 LIMIT 1
`

func (q *Queries) GetCurrentPlayerTurn(ctx context.Context, id int32) (int32, error) {
	row := q.db.QueryRow(ctx, getCurrentPlayerTurn, id)
	var currentplayerturn int32
	err := row.Scan(&currentplayerturn)
	return currentplayerturn, err
}

const getDeck = `-- name: GetDeck :one
SELECT deck.id, deck.cards FROM deck WHERE id=$1 LIMIT 1
`

func (q *Queries) GetDeck(ctx context.Context, id int32) (Deck, error) {
	row := q.db.QueryRow(ctx, getDeck, id)
	var i Deck
	err := row.Scan(&i.ID, &i.Cards)
	return i, err
}

const getGameCardsForMatch = `-- name: GetGameCardsForMatch :many
SELECT  m.id, m.match_id, m.cardid, m.origowner, m.currowner, m.state, c.id, c.value, c.suit, c.art
FROM matchcards m
JOIN cards c ON (m.cardId = c.id)
WHERE m.match_id = $1
`

type GetGameCardsForMatchRow struct {
	Matchcard Matchcard
	Card      Card
}

func (q *Queries) GetGameCardsForMatch(ctx context.Context, matchID int32) ([]GetGameCardsForMatchRow, error) {
	rows, err := q.db.Query(ctx, getGameCardsForMatch, matchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetGameCardsForMatchRow
	for rows.Next() {
		var i GetGameCardsForMatchRow
		if err := rows.Scan(
			&i.Matchcard.ID,
			&i.Matchcard.MatchID,
			&i.Matchcard.Cardid,
			&i.Matchcard.Origowner,
			&i.Matchcard.Currowner,
			&i.Matchcard.State,
			&i.Card.ID,
			&i.Card.Value,
			&i.Card.Suit,
			&i.Card.Art,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getMatchById = `-- name: GetMatchById :one
SELECT
    json_build_object(
        'id', id,
        'playerIds', playerIds,
        'creationDate', creationDate,
        'privateMatch', privateMatch,
        'eloRangeMin', eloRangeMin,
        'eloRangeMax', eloRangeMax,
        'deckid', deckid,
        'cutgamecardid', cutgamecardid,
        'currentplayerturn', currentplayerturn,
        'turnpasstimestamps', turnpasstimestamps,
        'art', art,
        'gameState', gameState,
        'players',
        (
            SELECT
                json_agg(
                    json_build_object(
                        'id', p.id,
                        'accountid', p.accountid,
                        'play', p.play,
                        'hand', p.hand,
                        'kitty', p.kitty,
                        'score', p.score,
                        'isready', p.isready,
                        'art', p.art
                    )
                )
            FROM player AS p
            WHERE p.Id = ANY(m.playerIds)
        )
    )
FROM match AS m
WHERE m.id = $1
LIMIT 1
`

func (q *Queries) GetMatchById(ctx context.Context, id int32) ([]byte, error) {
	row := q.db.QueryRow(ctx, getMatchById, id)
	var json_build_object []byte
	err := row.Scan(&json_build_object)
	return json_build_object, err
}

const getMatchByPlayerId = `-- name: GetMatchByPlayerId :one
SELECT
    json_build_object(
        'id', id,
        'playerIds', playerIds,
        'creationDate', creationDate,
        'privateMatch', privateMatch,
        'eloRangeMin', eloRangeMin,
        'eloRangeMax', eloRangeMax,
        'deckid', deckid,
        'cutgamecardid', cutgamecardid,
        'currentplayerturn', currentplayerturn,
        'turnpasstimestamps', turnpasstimestamps,
        'art', art,
        'gameState', gameState,
        'players',
        (
            SELECT
                json_agg(
                    json_build_object(
                        'id', p.id,
                        'accountid', p.accountid,
                        'play', p.play,
                        'hand', p.hand,
                        'kitty', p.kitty,
                        'score', p.score,
                        'isready', p.isready,
                        'art', p.art
                    )
                )
            FROM player AS p
            WHERE p.Id = ANY(m.playerIds)
        )
    )
FROM match as m 
WHERE $1::int=ANY(m.playerIds)
LIMIT 1
`

func (q *Queries) GetMatchByPlayerId(ctx context.Context, dollar_1 int32) ([]byte, error) {
	row := q.db.QueryRow(ctx, getMatchByPlayerId, dollar_1)
	var json_build_object []byte
	err := row.Scan(&json_build_object)
	return json_build_object, err
}

const getMatchCards = `-- name: GetMatchCards :many
SELECT matchcards.id, matchcards.match_id, matchcards.cardid, matchcards.origowner, matchcards.currowner, matchcards.state FROM matchcards NATURAL JOIN cards WHERE matchcards.id IN ($1::int[])
`

func (q *Queries) GetMatchCards(ctx context.Context, dollar_1 []int32) ([]Matchcard, error) {
	rows, err := q.db.Query(ctx, getMatchCards, dollar_1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Matchcard
	for rows.Next() {
		var i Matchcard
		if err := rows.Scan(
			&i.ID,
			&i.MatchID,
			&i.Cardid,
			&i.Origowner,
			&i.Currowner,
			&i.State,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getMatchIdForPlayerId = `-- name: GetMatchIdForPlayerId :one
SELECT id from match WHERE $1 = ANY(playerids) LIMIT 1
`

func (q *Queries) GetMatchIdForPlayerId(ctx context.Context, playerids []int32) (int32, error) {
	row := q.db.QueryRow(ctx, getMatchIdForPlayerId, playerids)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const getOpenMatches = `-- name: GetOpenMatches :many
SELECT
    json_build_object(
        'id', id,
        'playerIds', playerIds,
        'creationDate', creationDate,
        'privateMatch', privateMatch,
        'eloRangeMin', eloRangeMin,
        'eloRangeMax', eloRangeMax,
        'deckid', deckid,
        'cutgamecardid', cutgamecardid,
        'currentplayerturn', currentplayerturn,
        'turnpasstimestamps', turnpasstimestamps,
        'art', art,
        'gameState', gameState,
        'players',
        (
            SELECT
                json_agg(
                    json_build_object(
                        'id', p.id,
                        'accountid', p.accountid,
                        'play', p.play,
                        'hand', p.hand,
                        'kitty', p.kitty,
                        'score', p.score,
                        'isready', p.isready,
                        'art', p.art
                    )
                )
            FROM player AS p
            WHERE p.Id = ANY(m.playerIds)
        )
    )
FROM match AS m
`

func (q *Queries) GetOpenMatches(ctx context.Context) ([][]byte, error) {
	rows, err := q.db.Query(ctx, getOpenMatches)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items [][]byte
	for rows.Next() {
		var json_build_object []byte
		if err := rows.Scan(&json_build_object); err != nil {
			return nil, err
		}
		items = append(items, json_build_object)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPlayer = `-- name: GetPlayer :one
SELECT player.id, player.accountid, player.play, player.hand, player.kitty, player.score, player.isready, player.art FROM player WHERE id=$1 LIMIT 1
`

func (q *Queries) GetPlayer(ctx context.Context, id int32) (Player, error) {
	row := q.db.QueryRow(ctx, getPlayer, id)
	var i Player
	err := row.Scan(
		&i.ID,
		&i.Accountid,
		&i.Play,
		&i.Hand,
		&i.Kitty,
		&i.Score,
		&i.Isready,
		&i.Art,
	)
	return i, err
}

const passTurn = `-- name: PassTurn :exec
UPDATE match m
SET currentplayerturn = 
    (SELECT 
    CASE WHEN 
            array_position(playerids, currentplayerturn)=
            array_length(playerids,1)
        THEN 
            playerids[1]
        ELSE 
            playerids[array_position(playerids, currentplayerturn)+1]
        END
    FROM match m
    WHERE m.id = $1
    )            
WHERE m.id = $1
`

func (q *Queries) PassTurn(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, passTurn, id)
	return err
}

const removeCardsFromHand = `-- name: RemoveCardsFromHand :exec
UPDATE player SET hand = hand - $1 where id = $2
`

type RemoveCardsFromHandParams struct {
	Hand []int32
	ID   int32
}

func (q *Queries) RemoveCardsFromHand(ctx context.Context, arg RemoveCardsFromHandParams) error {
	_, err := q.db.Exec(ctx, removeCardsFromHand, arg.Hand, arg.ID)
	return err
}

const uodateGameState = `-- name: UodateGameState :exec
UPDATE match SET gameState= $1 WHERE id=$2
`

type UodateGameStateParams struct {
	Gamestate Gamestate
	ID        int32
}

func (q *Queries) UodateGameState(ctx context.Context, arg UodateGameStateParams) error {
	_, err := q.db.Exec(ctx, uodateGameState, arg.Gamestate, arg.ID)
	return err
}

const updateAccount = `-- name: UpdateAccount :exec
UPDATE match SET cutGameCardId = $2 where id=$1
`

type UpdateAccountParams struct {
	ID     int32
	Cardid int32
}

func (q *Queries) UpdateAccount(ctx context.Context, arg UpdateAccountParams) error {
	_, err := q.db.Exec(ctx, updateAccount, arg.ID, arg.Cardid)
	return err
}

const updateCardsPlayed = `-- name: UpdateCardsPlayed :exec
UPDATE player SET play = play + $1 where id = $2
`

type UpdateCardsPlayedParams struct {
	Play []int32
	ID   int32
}

func (q *Queries) UpdateCardsPlayed(ctx context.Context, arg UpdateCardsPlayedParams) error {
	_, err := q.db.Exec(ctx, updateCardsPlayed, arg.Play, arg.ID)
	return err
}

const updateKitty = `-- name: UpdateKitty :exec
UPDATE player SET kitty = kitty + $1 where id = $2
`

type UpdateKittyParams struct {
	Kitty []int32
	ID    int32
}

func (q *Queries) UpdateKitty(ctx context.Context, arg UpdateKittyParams) error {
	_, err := q.db.Exec(ctx, updateKitty, arg.Kitty, arg.ID)
	return err
}

const updateMatch = `-- name: UpdateMatch :exec
UPDATE match SET
	playerIds = $1,
	creationDate = $2,
	privateMatch = $3,
	eloRangeMin = $4,
	eloRangeMax = $5,
	deckId = $6,
	cutGameCardId = $7,
	currentPlayerTurn = $8,
	turnPassTimestamps = $9,
	gameState= $10,
	art = $11
WHERE id=$12
`

type UpdateMatchParams struct {
	Playerids          []int32
	Creationdate       pgtype.Timestamptz
	Privatematch       bool
	Elorangemin        int32
	Elorangemax        int32
	Deckid             int32
	Cutgamecardid      int32
	Currentplayerturn  int32
	Turnpasstimestamps []pgtype.Timestamptz
	Gamestate          Gamestate
	Art                string
	ID                 int32
}

func (q *Queries) UpdateMatch(ctx context.Context, arg UpdateMatchParams) error {
	_, err := q.db.Exec(ctx, updateMatch,
		arg.Playerids,
		arg.Creationdate,
		arg.Privatematch,
		arg.Elorangemin,
		arg.Elorangemax,
		arg.Deckid,
		arg.Cutgamecardid,
		arg.Currentplayerturn,
		arg.Turnpasstimestamps,
		arg.Gamestate,
		arg.Art,
		arg.ID,
	)
	return err
}

const updateMatchCut = `-- name: UpdateMatchCut :exec
UPDATE match SET cutGameCardId= $1 WHERE id=$2
`

type UpdateMatchCutParams struct {
	Cutgamecardid int32
	ID            int32
}

func (q *Queries) UpdateMatchCut(ctx context.Context, arg UpdateMatchCutParams) error {
	_, err := q.db.Exec(ctx, updateMatchCut, arg.Cutgamecardid, arg.ID)
	return err
}

const updateMatchState = `-- name: UpdateMatchState :exec
UPDATE match SET
	gameState= $1
WHERE id=$2
`

type UpdateMatchStateParams struct {
	Gamestate Gamestate
	ID        int32
}

func (q *Queries) UpdateMatchState(ctx context.Context, arg UpdateMatchStateParams) error {
	_, err := q.db.Exec(ctx, updateMatchState, arg.Gamestate, arg.ID)
	return err
}

const updateMatchWithDeckId = `-- name: UpdateMatchWithDeckId :exec
UPDATE match SET deckid = $1 where id = $2
`

type UpdateMatchWithDeckIdParams struct {
	Deckid int32
	ID     int32
}

func (q *Queries) UpdateMatchWithDeckId(ctx context.Context, arg UpdateMatchWithDeckIdParams) error {
	_, err := q.db.Exec(ctx, updateMatchWithDeckId, arg.Deckid, arg.ID)
	return err
}

const updatePlayer = `-- name: UpdatePlayer :one
UPDATE player SET 
		hand = $1, 
		play = $2, 
		kitty = $3, 
		score = $4, 
		isReady = $5,
		art = $6 
	WHERE 
		id = $7
    RETURNING id, accountid, play, hand, kitty, score, isready, art
`

type UpdatePlayerParams struct {
	Hand    []int32
	Play    []int32
	Kitty   []int32
	Score   int32
	Isready bool
	Art     string
	ID      int32
}

func (q *Queries) UpdatePlayer(ctx context.Context, arg UpdatePlayerParams) (Player, error) {
	row := q.db.QueryRow(ctx, updatePlayer,
		arg.Hand,
		arg.Play,
		arg.Kitty,
		arg.Score,
		arg.Isready,
		arg.Art,
		arg.ID,
	)
	var i Player
	err := row.Scan(
		&i.ID,
		&i.Accountid,
		&i.Play,
		&i.Hand,
		&i.Kitty,
		&i.Score,
		&i.Isready,
		&i.Art,
	)
	return i, err
}

const updatePlayerReady = `-- name: UpdatePlayerReady :exec
UPDATE player SET isReady = $1 WHERE id = $2
`

type UpdatePlayerReadyParams struct {
	Isready bool
	ID      int32
}

func (q *Queries) UpdatePlayerReady(ctx context.Context, arg UpdatePlayerReadyParams) error {
	_, err := q.db.Exec(ctx, updatePlayerReady, arg.Isready, arg.ID)
	return err
}

const updatePlayersInMatch = `-- name: UpdatePlayersInMatch :exec

UPDATE match SET playerIds=ARRAY_APPEND(playerIds, $1) WHERE id=$2
`

type UpdatePlayersInMatchParams struct {
	ArrayAppend interface{}
	ID          int32
}

// -- name: StartMatch :exec
// UPDATE match SET playerIds=ARRAY_APPEND(playerIds, $1), deckid=$2 WHERE id=$3;
func (q *Queries) UpdatePlayersInMatch(ctx context.Context, arg UpdatePlayersInMatchParams) error {
	_, err := q.db.Exec(ctx, updatePlayersInMatch, arg.ArrayAppend, arg.ID)
	return err
}
