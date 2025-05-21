-- name: CreateMatch :one
INSERT INTO match(
    privateMatch,
    eloRangeMin,
    eloRangeMax,
    cutGameCardId,
    turnPassTimestamps,
    dealerId,
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
  $8)
RETURNING *;

-- name: UpdateMatchCut :exec
UPDATE match SET cutGameCardId= $1 WHERE id=$2;

-- name: UpdateMatchState :one
UPDATE match SET gameState= $1 WHERE id=$2 RETURNING *;

-- name: UpdateMatch :exec
UPDATE match SET
	creationDate = $1,
	privateMatch = $2,
	eloRangeMin = $3,
	eloRangeMax = $4,
	cutGameCardId = $5,
  dealerId = $6,
	currentPlayerTurn = $7,
	turnPassTimestamps = $8,
	gameState= $9,
	art = $10
WHERE id=$11;

-- name: UpateMatchCurrentPlayerTurn :exec
UPDATE match SET currentplayerturn = $1 WHERE id = $2;

-- name: UpdateDealerForMatch :exec
UPDATE match SET dealerid = $1 WHERE id = $2;

-- name: ResetDeckForMatchId :exec
UPDATE matchcard m SET state = 'Deck', origowner = null, currowner = null FROM matchcard
LEFT JOIN
    deck ON matchcard.deckid=deck.id
LEFT JOIN
    match ON deck.matchId=match.id
WHERE
    match.id = $1;

-- name: GetMatchById :one
SELECT
    json_build_object(
        'id', id,
        'creationDate', creationDate,
        'privateMatch', privateMatch,
        'eloRangeMin', eloRangeMin,
        'eloRangeMax', eloRangeMax,
        'deckid', deckid,
        'cutgamecardid', cutgamecardid,
        'dealerid', dealerid,
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
LIMIT 1;

-- name: GetOpenMatches :many
SELECT * FROM match WHERE gameState=$1;

-- name: GetMatchCurrentPlayerTurn :one
SELECT currentplayerturn FROM match WHERE id = $1 LIMIT 1;

-- name: GetNextPlayerInTurnOrder :one
SELECT 
    player.*,
    match_player.turnorder
FROM
    player
LEFT JOIN
    match_player ON player.id=match_player.playerid
WHERE
    match_player.matchid = $1 AND match_player.turnorder = $2 + 1;

-- name: GetPlayersForMatchId :many
SELECT
    json_build_object(
        'id', p.id,
        'accountid', p.accountid,
        'score', p.score,
        'isready', p.isready,
        'art', p.art,
        'hand',
        (
            SELECT
                json_agg(
                    json_build_object(
                        'id', m.id,
                        'cardid', m.cardid,
                        'origowner', m.origowner,
                        'currowner', m.currowner,
                        'state', m.state
                    )
                )
            FROM matchcard AS m
            WHERE m.currowner = p.id AND m.state = 'Hand'
        ),
        'kitty',
        (
            SELECT
                json_agg(
                    json_build_object(
                        'id', m.id,
                        'cardid', m.cardid,
                        'origowner', m.origowner,
                        'currowner', m.currowner,
                        'state', m.state
                    )
                )
            FROM matchcard AS m
            WHERE m.currowner = p.id AND m.state = 'Kitty'
        ),
        'play',
        (
            SELECT
                json_agg(
                    json_build_object(
                        'id', m.id,
                        'cardid', m.cardid,
                        'origowner', m.origowner,
                        'currowner', m.currowner,
                        'state', m.state
                    )
                )
            FROM matchcard AS m
            WHERE m.state = 'Play'
        )
    )
FROM player as p
LEFT JOIN
    match_player ON p.id=match_player.playerid
WHERE
    match_player.matchid = $1;

-- name: GetMatchStateById :one
SELECT 
    match.gameState
FROM
    match
WHERE
    match.id = $1;

-- name: GetDeckForMatchId :many
SELECT 
    deck.*,
    matchcard.*,
    card.*
FROM
    matchcard
LEFT JOIN
    deck ON matchcard.deckid=deck.id
LEFT JOIN
    match ON deck.matchid=match.id
left join 
	card on card.id=matchcard.cardid
WHERE
    match.id = $1;

