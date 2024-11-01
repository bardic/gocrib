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
    card ON deck_matchcard.cardId=card.id
WHERE
    deck.id IN ($1);

-- name: UpdateCardsPlayed :exec
UPDATE player SET play = play + $1 where id = $2;

-- name: RemoveCardsFromHand :exec
UPDATE player SET hand = hand - $1 where id = $2;

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
		hand = $1, 
		play = $2, 
		kitty = $3, 
		score = $4, 
		isReady = $5,
		art = $6 
	WHERE 
		id = $7
    RETURNING *;

-- name: CreateMatch :one
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

-- name: UpdateMatchWithDeckId :exec
UPDATE match SET deckid = $1 where id = $2;

-- name: CreateMatchCards :one
INSERT INTO matchcard (cardid, origowner, currowner, state) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: PassTurn :exec
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
WHERE m.id = $1;

-- name: JoinMatch :exec
INSERT INTO 
    match_player (matchid, playerid)     
VALUES 
    ($1, $2);