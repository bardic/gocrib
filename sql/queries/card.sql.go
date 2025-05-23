// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: card.sql

package queries

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type CreateMatchCardParams struct {
	Cardid    int
	Origowner int
	Currowner int
	State     string
	Deckid    int
}

const getCards = `-- name: GetCards :many
SELECT id, value, suit, art FROM card
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

const getCardsForMatchId = `-- name: GetCardsForMatchId :many
SELECT 
    deck.id, deck.cutmatchcardid, deck.matchid,
    matchcard.id, matchcard.cardid, matchcard.origowner, matchcard.currowner, matchcard.state, matchcard.deckid,
    card.id, card.value, card.suit, card.art
FROM
    matchcard
LEFT JOIN
    deck ON matchcard.deckid=deck.id
LEFT JOIN
    match ON deck.matchId=match.id
left join 
	card on card.id=matchcard.cardid
WHERE
    match.id = $1
`

type GetCardsForMatchIdRow struct {
	ID             int
	Cutmatchcardid int
	Matchid        int
	ID_2           int
	Cardid         int
	Origowner      int
	Currowner      int
	State          string
	Deckid         int
	ID_3           int
	Value          NullCardvalue
	Suit           NullCardsuit
	Art            pgtype.Text
}

func (q *Queries) GetCardsForMatchId(ctx context.Context, id int) ([]GetCardsForMatchIdRow, error) {
	rows, err := q.db.Query(ctx, getCardsForMatchId, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetCardsForMatchIdRow
	for rows.Next() {
		var i GetCardsForMatchIdRow
		if err := rows.Scan(
			&i.ID,
			&i.Cutmatchcardid,
			&i.Matchid,
			&i.ID_2,
			&i.Cardid,
			&i.Origowner,
			&i.Currowner,
			&i.State,
			&i.Deckid,
			&i.ID_3,
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

const getCardsForMatchIdAndState = `-- name: GetCardsForMatchIdAndState :many
SELECT 
    deck.id, deck.cutmatchcardid, deck.matchid,
    matchcard.id, matchcard.cardid, matchcard.origowner, matchcard.currowner, matchcard.state, matchcard.deckid,
    card.id, card.value, card.suit, card.art
FROM
    matchcard
LEFT JOIN
    deck ON matchcard.deckid=deck.id
LEFT JOIN
    match ON deck.matchId=match.id
left join 
	card on card.id=matchcard.cardid
WHERE
    match.id = $1 AND matchcard.state = $2
`

type GetCardsForMatchIdAndStateParams struct {
	ID    int
	State string
}

type GetCardsForMatchIdAndStateRow struct {
	ID             int
	Cutmatchcardid int
	Matchid        int
	ID_2           int
	Cardid         int
	Origowner      int
	Currowner      int
	State          string
	Deckid         int
	ID_3           int
	Value          NullCardvalue
	Suit           NullCardsuit
	Art            pgtype.Text
}

func (q *Queries) GetCardsForMatchIdAndState(ctx context.Context, arg GetCardsForMatchIdAndStateParams) ([]GetCardsForMatchIdAndStateRow, error) {
	rows, err := q.db.Query(ctx, getCardsForMatchIdAndState, arg.ID, arg.State)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetCardsForMatchIdAndStateRow
	for rows.Next() {
		var i GetCardsForMatchIdAndStateRow
		if err := rows.Scan(
			&i.ID,
			&i.Cutmatchcardid,
			&i.Matchid,
			&i.ID_2,
			&i.Cardid,
			&i.Origowner,
			&i.Currowner,
			&i.State,
			&i.Deckid,
			&i.ID_3,
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

const getCardsForPlayerIdFromDeckId = `-- name: GetCardsForPlayerIdFromDeckId :many
SELECT
    card.id, card.value, card.suit, card.art,
    matchcard.id, matchcard.cardid, matchcard.origowner, matchcard.currowner, matchcard.state, matchcard.deckid 
FROM
    card
LEFT JOIN
    matchcard ON card.id=matchcard.cardid
WHERE
    matchcard.deckid = $1 AND matchcard.origowner = $2
`

type GetCardsForPlayerIdFromDeckIdParams struct {
	Deckid    int
	Origowner int
}

type GetCardsForPlayerIdFromDeckIdRow struct {
	ID        int
	Value     Cardvalue
	Suit      Cardsuit
	Art       string
	ID_2      int
	Cardid    int
	Origowner int
	Currowner int
	State     pgtype.Text
	Deckid    int
}

func (q *Queries) GetCardsForPlayerIdFromDeckId(ctx context.Context, arg GetCardsForPlayerIdFromDeckIdParams) ([]GetCardsForPlayerIdFromDeckIdRow, error) {
	rows, err := q.db.Query(ctx, getCardsForPlayerIdFromDeckId, arg.Deckid, arg.Origowner)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetCardsForPlayerIdFromDeckIdRow
	for rows.Next() {
		var i GetCardsForPlayerIdFromDeckIdRow
		if err := rows.Scan(
			&i.ID,
			&i.Value,
			&i.Suit,
			&i.Art,
			&i.ID_2,
			&i.Cardid,
			&i.Origowner,
			&i.Currowner,
			&i.State,
			&i.Deckid,
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

const updateMatchCardState = `-- name: UpdateMatchCardState :exec
UPDATE matchcard SET state = $1, origowner = $2, currowner = $3 WHERE id = $4
`

type UpdateMatchCardStateParams struct {
	State     string
	Origowner int
	Currowner int
	ID        int
}

func (q *Queries) UpdateMatchCardState(ctx context.Context, arg UpdateMatchCardStateParams) error {
	_, err := q.db.Exec(ctx, updateMatchCardState,
		arg.State,
		arg.Origowner,
		arg.Currowner,
		arg.ID,
	)
	return err
}
