-- name: GetMatchById :one
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
LIMIT 1;


-- name: GetMatchByPlayerId :one
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
LIMIT 1;

-- name: GetOpenMatches :many
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
FROM match AS m;

-- name: GetAccount :one
SELECT accounts.* FROM accounts WHERE id = $1 LIMIT 1;

-- name: UpdateAccount :exec
UPDATE match SET cutGameCardId = @cardId where id=$1;

-- name: GetCards :many
SELECT cards.* FROM cards;

-- name: CreateDeck :one
INSERT INTO deck(cards) VALUES ($1) RETURNING *;

-- name: UodateGameState :exec
UPDATE match SET gameState= $1 WHERE id=$2;

-- name: UpdateMatchCut :exec
UPDATE match SET cutGameCardId= $1 WHERE id=$2;

-- name: UpdateMatch :exec
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
WHERE id=$12;

-- name: StartMatch :exec
UPDATE match SET playerIds=ARRAY_APPEND(playerIds, $1), deckid=$2 WHERE id=$3;

-- name: GetDeck :one
SELECT deck.* FROM deck WHERE id=$1 LIMIT 1;

-- name: GetMatchCards :many
SELECT matchcards.* FROM matchcards NATURAL JOIN cards WHERE matchcards.id IN ($1::int[]);

-- name: UpdateCardsPlayed :exec
UPDATE player SET play = play + $1 where id = $2;

-- name: RemoveCardsFromHand :exec
UPDATE player SET hand = hand - $1 where id = $2;

-- name: GetMatchIdForPlayerId :one 
SELECT id from match WHERE $1 = ANY(playerids) LIMIT 1;

-- name: GetPlayer :one
SELECT player.* FROM player WHERE id=$1 LIMIT 1;

-- name: UpdatePlayer :exec
UPDATE player SET 
		hand = $1, 
		play = $2, 
		kitty = $3, 
		score = $4, 
		isReady = $5,
		art = $6 
	where 
		id = $7;

-- name: CreateMatch :one
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
			RETURNING *;

-- name: CreatePlayer :one
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
		RETURNING *;

-- name: GetCurrentPlayerTurn :one
SELECT currentplayerturn FROM match WHERE id = $1 LIMIT 1;

-- name: UpdatePlayerReady :exec
UPDATE player SET isReady = $1 WHERE id = $2;

-- name: UpdateKitty :exec
UPDATE player SET kitty = kitty + $1 where id = $2;