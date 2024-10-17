-- name: GetMatch :one
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
    1;

-- name: GetAccount :one
SELECT name FROM accounts WHERE id=$1;