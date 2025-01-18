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


-- name: GetMatchByPlayerId :one
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
LIMIT 1;

-- name: GetOpenMatches :many
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
    match_player as mp ON m.id=mp.matchid;

-- name: GetAccount :one
SELECT account.* FROM account WHERE id = $1 LIMIT 1;

-- name: UpdateAccount :exec
UPDATE match SET cutGameCardId = @cardId where id=$1;

-- name: GetCards :many
SELECT card.* FROM card;

-- name: CreateDeck :one
INSERT INTO deck(cutmatchcardid) VALUES (null) RETURNING *;

-- name: UpdateGameState :one
UPDATE match SET gameState= $1 WHERE id=$2 RETURNING *;

-- name: UpdateMatchCut :exec
UPDATE match SET cutGameCardId= $1 WHERE id=$2;

-- name: UpdateMatch :exec
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
WHERE id=$11;

-- name: UpdateMatchState :exec
UPDATE match SET
	gameState= $1
WHERE id=$2;

-- name: GetDeckForMatchId :one
SELECT deck.* FROM deck
LEFT JOIN
    match ON deck.id=match.deckid
 WHERE match.id=$1 LIMIT 1;

-- name: GetMatchCards :many
SELECT 
    sqlc.embed(deck_matchcard),
    sqlc.embed(deck),
    sqlc.embed(matchcard),
    sqlc.embed(card)
FROM 
    deck_matchcard
LEFT JOIN
    matchcard ON deck_matchcard.matchcardid=matchcard.id
LEFT JOIN
    deck ON deck_matchcard.deckid=deck.id
LEFT JOIN
    card ON deck_matchcard.matchcardId=card.id
WHERE
    deck.id IN ($1);

-- name: GetMatchIdForPlayerId :one 
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

-- name: GetPlayer :one
SELECT player.* FROM player WHERE id=$1 LIMIT 1;

-- name: UpdatePlayer :one
UPDATE player SET 
		score = $1, 
		isReady = $2,
		art = $3 
	WHERE 
		id = $4
    RETURNING *;

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

-- name: CreatePlayer :one
INSERT INTO player (
			accountid,
			score,
			isready,
			art
		) VALUES (
			$1,
			$2,
			$3,
			$4
		)
		RETURNING *;

-- name: GetCurrentPlayerTurn :one
SELECT currentplayerturn FROM match WHERE id = $1 LIMIT 1;

-- name: UpdatePlayerReady :exec
UPDATE player SET isReady = $1 WHERE id = $2;

-- name: UpdateMatchWithDeckId :exec
UPDATE match SET deckid = $1 where id = $2;

-- name: CreateMatchCards :one
INSERT INTO matchcard (cardid, origowner, currowner, state) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: PassTurn :exec
UPDATE match m
SET currentplayerturn = 
    (
    SELECT playerId FROM match_player p
 where p.playerId !=$2 and p.matchId =$1
    )    
WHERE m.id = $1;

-- name: JoinMatch :exec
INSERT INTO 
    match_player (matchid, playerid)     
VALUES 
    ($1, $2);

-- name: InsertDeckMatchCard :exec
INSERT INTO deck_matchcard (deckid, matchcardid) VALUES ($1, $2);

-- name: UpdateCurrentPlayerTurn :exec
UPDATE match SET currentplayerturn = $1 WHERE id = $2;

-- -- name: RemoveCardFromDeck :exec
-- DELETE FROM deck_matchcard WHERE deckid = $1 AND matchcardid = $2;

-- name: GetPlayersInMatch :many
SELECT 
    player.*
FROM
    player
LEFT JOIN
    match_player ON player.id=match_player.playerid
WHERE
    match_player.matchid = $1;  

-- name: UpdateMatchCardState :exec
UPDATE matchcard SET state = $1, origowner = $2, currowner = $3 WHERE id = $4;