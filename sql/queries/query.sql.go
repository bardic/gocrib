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
INSERT INTO deck(cutmatchcardid) VALUES (null) RETURNING id, cutmatchcardid
`

func (q *Queries) CreateDeck(ctx context.Context) (Deck, error) {
	row := q.db.QueryRow(ctx, createDeck)
	var i Deck
	err := row.Scan(&i.ID, &i.Cutmatchcardid)
	return i, err
}

const createMatch = `-- name: CreateMatch :one
INSERT INTO match(
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
				$9)
			RETURNING id, creationdate, privatematch, elorangemin, elorangemax, deckid, cutgamecardid, currentplayerturn, turnpasstimestamps, gamestate, art
`

type CreateMatchParams struct {
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
INSERT INTO matchcard (cardid, origowner, currowner, state) VALUES ($1, $2, $3, $4) RETURNING id, cardid, origowner, currowner, state
`

type CreateMatchCardsParams struct {
	Cardid    int32
	Origowner pgtype.Int4
	Currowner pgtype.Int4
	State     Cardstate
}

func (q *Queries) CreateMatchCards(ctx context.Context, arg CreateMatchCardsParams) (Matchcard, error) {
	row := q.db.QueryRow(ctx, createMatchCards,
		arg.Cardid,
		arg.Origowner,
		arg.Currowner,
		arg.State,
	)
	var i Matchcard
	err := row.Scan(
		&i.ID,
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
SELECT account.id, account.name FROM account WHERE id = $1 LIMIT 1
`

func (q *Queries) GetAccount(ctx context.Context, id int32) (Account, error) {
	row := q.db.QueryRow(ctx, getAccount, id)
	var i Account
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const getCards = `-- name: GetCards :many
SELECT card.id, card.value, card.suit, card.art FROM card
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

const getDeckForMatchId = `-- name: GetDeckForMatchId :one
SELECT deck.id, deck.cutmatchcardid FROM deck
LEFT JOIN
    match ON deck.id=match.deckid
 WHERE match.id=$1 LIMIT 1
`

func (q *Queries) GetDeckForMatchId(ctx context.Context, id int32) (Deck, error) {
	row := q.db.QueryRow(ctx, getDeckForMatchId, id)
	var i Deck
	err := row.Scan(&i.ID, &i.Cutmatchcardid)
	return i, err
}

const getMatchById = `-- name: GetMatchById :one
SELECT
    json_build_object(
        'id', id,
        'creationDate', creationDate,
        'privateMatch', privateMatch,
        'eloRangeMin', eloRangeMin,
        'eloRangeMax', eloRangeMax,
        'deckid', deckid,
        'cutgamecardid', cutgamecardid,
        'currentplayerturn', currentplayerturn,
        'turnpasstimestamps', turnpasstimestamps,
        'gameState', gameState,
        'art', art,
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
            LEFT JOIN
                match_player as mp ON p.id=mp.playerid and mp.matchid=m.id
            WHERE p.Id = mp.playerId
        )
    )
FROM match AS m
LEFT JOIN
    match_player as mp ON m.id=mp.matchid
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
        'creationDate', creationDate,
        'privateMatch', privateMatch,
        'eloRangeMin', eloRangeMin,
        'eloRangeMax', eloRangeMax,
        'deckid', deckid,
        'cutgamecardid', cutgamecardid,
        'currentplayerturn', currentplayerturn,
        'turnpasstimestamps', turnpasstimestamps,
        'gameState', gameState,
        'art', art,
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
LIMIT 1
`

func (q *Queries) GetMatchByPlayerId(ctx context.Context) ([]byte, error) {
	row := q.db.QueryRow(ctx, getMatchByPlayerId)
	var json_build_object []byte
	err := row.Scan(&json_build_object)
	return json_build_object, err
}

const getMatchCards = `-- name: GetMatchCards :many
SELECT 
    deck_matchcard.deckid, deck_matchcard.matchcardid,
    deck.id, deck.cutmatchcardid,
    matchcard.id, matchcard.cardid, matchcard.origowner, matchcard.currowner, matchcard.state,
    card.id, card.value, card.suit, card.art
FROM 
    deck_matchcard
LEFT JOIN
    matchcard ON deck_matchcard.matchcardid=matchcard.id
LEFT JOIN
    deck ON deck_matchcard.deckid=deck.id
LEFT JOIN
    card ON deck_matchcard.matchcardId=card.id
WHERE
    deck.id IN ($1)
`

type GetMatchCardsRow struct {
	DeckMatchcard DeckMatchcard
	Deck          Deck
	Matchcard     Matchcard
	Card          Card
}

func (q *Queries) GetMatchCards(ctx context.Context, id int32) ([]GetMatchCardsRow, error) {
	rows, err := q.db.Query(ctx, getMatchCards, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetMatchCardsRow
	for rows.Next() {
		var i GetMatchCardsRow
		if err := rows.Scan(
			&i.DeckMatchcard.Deckid,
			&i.DeckMatchcard.Matchcardid,
			&i.Deck.ID,
			&i.Deck.Cutmatchcardid,
			&i.Matchcard.ID,
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

const getMatchIdForPlayerId = `-- name: GetMatchIdForPlayerId :one
SELECT 
    match_player.matchid, match_player.playerid,
    match.id, match.creationdate, match.privatematch, match.elorangemin, match.elorangemax, match.deckid, match.cutgamecardid, match.currentplayerturn, match.turnpasstimestamps, match.gamestate, match.art,
    player.id, player.accountid, player.play, player.hand, player.kitty, player.score, player.isready, player.art
FROM 
    match_player
INNER JOIN
    match ON match_player.matchid=match.id
LEFT JOIN
    player ON match_player.playerid=player.id
WHERE $1 = match_player.playerId LIMIT 1
`

type GetMatchIdForPlayerIdRow struct {
	Matchid            int32
	Playerid           int32
	ID                 int32
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
	ID_2               pgtype.Int4
	Accountid          pgtype.Int4
	Play               []int32
	Hand               []int32
	Kitty              []int32
	Score              pgtype.Int4
	Isready            pgtype.Bool
	Art_2              pgtype.Text
}

func (q *Queries) GetMatchIdForPlayerId(ctx context.Context, playerid int32) (GetMatchIdForPlayerIdRow, error) {
	row := q.db.QueryRow(ctx, getMatchIdForPlayerId, playerid)
	var i GetMatchIdForPlayerIdRow
	err := row.Scan(
		&i.Matchid,
		&i.Playerid,
		&i.ID,
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
		&i.ID_2,
		&i.Accountid,
		&i.Play,
		&i.Hand,
		&i.Kitty,
		&i.Score,
		&i.Isready,
		&i.Art_2,
	)
	return i, err
}

const getOpenMatches = `-- name: GetOpenMatches :many
SELECT
    json_build_object(
        'id', id,
        'creationDate', creationDate,
        'privateMatch', privateMatch,
        'eloRangeMin', eloRangeMin,
        'eloRangeMax', eloRangeMax,
        'deckid', deckid,
        'cutgamecardid', cutgamecardid,
        'currentplayerturn', currentplayerturn,
        'turnpasstimestamps', turnpasstimestamps,
        'gameState', gameState,
        'art', art,
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
            LEFT JOIN
                match_player as mp ON p.id=mp.playerid
            WHERE p.Id = mp.playerId
        )
    )
FROM match AS m
LEFT JOIN
    match_player as mp ON m.id=mp.matchid
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

const insertDeckMatchCard = `-- name: InsertDeckMatchCard :exec
INSERT INTO deck_matchcard (deckid, matchcardid) VALUES ($1, $2)
`

type InsertDeckMatchCardParams struct {
	Deckid      int32
	Matchcardid int32
}

func (q *Queries) InsertDeckMatchCard(ctx context.Context, arg InsertDeckMatchCardParams) error {
	_, err := q.db.Exec(ctx, insertDeckMatchCard, arg.Deckid, arg.Matchcardid)
	return err
}

const joinMatch = `-- name: JoinMatch :exec
INSERT INTO 
    match_player (matchid, playerid)     
VALUES 
    ($1, $2)
`

type JoinMatchParams struct {
	Matchid  int32
	Playerid int32
}

func (q *Queries) JoinMatch(ctx context.Context, arg JoinMatchParams) error {
	_, err := q.db.Exec(ctx, joinMatch, arg.Matchid, arg.Playerid)
	return err
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

const updateGameState = `-- name: UpdateGameState :one
UPDATE match SET gameState= $1 WHERE id=$2 RETURNING id, creationdate, privatematch, elorangemin, elorangemax, deckid, cutgamecardid, currentplayerturn, turnpasstimestamps, gamestate, art
`

type UpdateGameStateParams struct {
	Gamestate Gamestate
	ID        int32
}

func (q *Queries) UpdateGameState(ctx context.Context, arg UpdateGameStateParams) (Match, error) {
	row := q.db.QueryRow(ctx, updateGameState, arg.Gamestate, arg.ID)
	var i Match
	err := row.Scan(
		&i.ID,
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
	creationDate = $1,
	privateMatch = $2,
	eloRangeMin = $3,
	eloRangeMax = $4,
	deckId = $5,
	cutGameCardId = $6,
	currentPlayerTurn = $7,
	turnPassTimestamps = $8,
	gameState= $9,
	art = $10
WHERE id=$11
`

type UpdateMatchParams struct {
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
