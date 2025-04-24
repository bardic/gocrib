-- name: CreateMatch :one
INSERT INTO match(
    privateMatch,
    eloRangeMin,
    eloRangeMax,
    deckId,
    cutGameCardId,
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
	deckId = $5,
	cutGameCardId = $6,
    dealerId = $7,
	currentPlayerTurn = $8,
	turnPassTimestamps = $9,
	gameState= $10,
	art = $11
WHERE id=$12;

-- name: UpateMatchCurrentPlayerTurn :exec
UPDATE match SET currentplayerturn = $1 WHERE id = $2;

-- name: UpdateMatchWithDeckId :exec
UPDATE match SET deckid = $1 where id = $2;

-- name: UpdateDealerForMatch :exec
UPDATE match SET dealerid = $1 WHERE id = $2;

-- name: ResetDeckForMatchId :exec
UPDATE matchcard m SET state = 'Deck', origowner = null, currowner = null FROM matchcard
LEFT JOIN 
    deck_matchcard ON matchcard.id=deck_matchcard.matchcardid
LEFT JOIN
    deck ON deck_matchcard.deckid=deck.id
LEFT JOIN
    match ON deck.id=match.deckid
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

-- name: GetMatchForPlayerId :one 
SELECT 
    match_player.*,
    match.*,
    player.*
FROM 
    match_player
INNER JOIN
    match ON match_player.matchid=match.id
LEFT JOIN
    player ON match_player.playerid=player.id
WHERE $1 = match_player.playerId LIMIT 1;

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

-- name: GetMatchDealer :one  
SELECT 
    player.*
FROM
    player
LEFT JOIN
    match ON player.id=match.dealerid
WHERE
    match.id = $1;

-- TODO this takes a second param. PRetty sure we can get that from the match object
-- name: GetMatchNextPlayerForMatchId :one
SELECT 
    player.*,
    match_player.turnorder
FROM
    player
LEFT JOIN
    match_player ON player.id=match_player.playerid
WHERE
    match_player.matchid = $1 AND match_player.turnorder = $2 + 1;
