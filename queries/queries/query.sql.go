// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package queries

import (
	"context"
)

const getAccount = `-- name: GetAccount :one
SELECT name FROM accounts WHERE id=$1
`

func (q *Queries) GetAccount(ctx context.Context, id int32) (string, error) {
	row := q.db.QueryRow(ctx, getAccount, id)
	var name string
	err := row.Scan(&name)
	return name, err
}

const getMatch = `-- name: GetMatch :one
SELECT
    json_build_object(
        'id',
        id,
        'playerIds',
        playerIds,
        'creationDate',
        creationDate,
        'privateMatch',
        privateMatch,
        'eloRangeMin',
        eloRangeMin,
        'eloRangeMax',
        eloRangeMax,
        'deckid',
        deckid,
        'cutgamecardid',
        cutgamecardid,
        'currentplayerturn',
        currentplayerturn,
        'turnpasstimestamps',
        turnpasstimestamps,
        'art',
        art,
        'gameState',
        gameState,
        'players',
        (
            SELECT
                json_agg(
                    json_build_object(
                        'id',
                        p.id,
                        'accountid',
                        p.accountid,
                        'play',
                        p.play,
                        'hand',
                        p.hand,
                        'kitty',
                        p.kitty,
                        'score',
                        p.score,
                        'isready',
                        p.isready,
                        'art',
                        p.art
                    )
                )
            FROM
                player as p
            WHERE
                p.Id = ANY(m.playerIds)
        )
    )
FROM
    match as m
WHERE
    m.id = $1
LIMIT
    1
`

func (q *Queries) GetMatch(ctx context.Context, id int32) ([]byte, error) {
	row := q.db.QueryRow(ctx, getMatch, id)
	var json_build_object []byte
	err := row.Scan(&json_build_object)
	return json_build_object, err
}