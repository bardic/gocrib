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
        'dealerid', dealerid,
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
    dealerId = $7,
	currentPlayerTurn = $8,
	turnPassTimestamps = $9,
	gameState= $10,
	art = $11
WHERE id=$12;

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
    deck_matchcard.*, 
    deck.*,
    matchcard.*,
    card.*
FROM 
    deck_matchcard
LEFT JOIN
    matchcard ON deck_matchcard.matchcardid=matchcard.id
LEFT JOIN
    deck ON deck_matchcard.deckid=deck.id
LEFT JOIN
    card ON deck_matchcard.matchcardId=card.id
WHERE
    deck.id = $1;

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

-- name: GetMatchPlayerOrdered :many
SELECT 
    player.*,
    match_player.turnOrder
FROM
    player
LEFT JOIN
    match_player ON player.id=match_player.playerid
WHERE
    match_player.matchid = $1
ORDER BY
    match_player.turnOrder ASC;

-- name: GetCardsForPlayerAndDeck :many
SELECT
    card.*,
    matchcard.state 
FROM
    card
LEFT JOIN
    matchcard ON card.id=matchcard.cardid
LEFT JOIN
    deck_matchcard ON matchcard.id=deck_matchcard.matchcardid
WHERE
    deck_matchcard.deckid = $1 AND matchcard.origowner = $2;

-- name: GetDealerForMatchId :one  
SELECT 
    player.*
FROM
    player
LEFT JOIN
    match ON player.id=match.dealerid
WHERE
    match.id = $1;

-- name: UpdateDealerForMatch :exec
UPDATE match SET dealerid = $1 WHERE id = $2;

-- name: GetPlayerById :one
SELECT 	
	p.*,
	array(	
		(select m.* from 
		matchcard m where
		m.currowner=1 and m.state='Hand'
		)) as hand,
	array(	
		(select m.* from 
		matchcard m where
		m.currowner=1 and m.state='Play'
		)) as board,
	array(	
		(select m.* from 
		matchcard m where
		m.currowner=$1 and m.state='Kitty'
		)) as kitty
FROM player as p
WHERE p.id=$1;

-- name: GetMarchCardsByType :many

SELECT 
    matchcard.*
FROM
    matchcard
LEFT JOIN
    deck_matchcard ON matchcard.id=deck_matchcard.matchcardid
LEFT JOIN
    deck ON deck_matchcard.deckid=deck.id
LEFT JOIN
    match ON deck.id=match.deckid
WHERE
    match.id = $1 AND matchcard.state = $2;

-- name: GetMatchCardsByTypeAndDeckId :many
SELECT 
    deck_matchcard.*, 
    deck.*,
    matchcard.*,
    card.*
FROM 
    deck_matchcard
LEFT JOIN
    matchcard ON deck_matchcard.matchcardid=matchcard.id
LEFT JOIN
    deck ON deck_matchcard.deckid=deck.id
LEFT JOIN
    card ON deck_matchcard.matchcardId=card.id
WHERE
    deck.id = $1 AND matchcard.state = $2;

-- name: GetMatchCardsByPlayerIdAndDeckId :many
SELECT 
    deck_matchcard.*, 
    deck.*,
    matchcard.*,
    card.*
FROM 
    deck_matchcard
LEFT JOIN
    matchcard ON deck_matchcard.matchcardid=matchcard.id
LEFT JOIN
    deck ON deck_matchcard.deckid=deck.id
LEFT JOIN
    card ON deck_matchcard.matchcardId=card.id
WHERE
     deck.id = $1 AND (matchcard.currowner = $2 OR matchcard.currowner IS NULL);

-- name: ResetDeckState :exec
UPDATE matchcard m SET state = 'Deck', origowner = null, currowner = null FROM matchcard
LEFT JOIN 
    deck_matchcard ON matchcard.id=deck_matchcard.matchcardid
LEFT JOIN
    deck ON deck_matchcard.deckid=deck.id
LEFT JOIN
    match ON deck.id=match.deckid
WHERE
    match.id = $1;


-- name: GetMatchPlayersByMatchId :many
SELECT 
    player.*,
    match_player.turnorder 
FROM
    player
LEFT JOIN   
    match_player ON player.id=match_player.playerid
WHERE
    match_player.matchid = $1;
    
-- name: UpdatePlayerTurnOrder :exec

UPDATE match_player SET turnorder = $1 WHERE matchid = $2 AND playerid = $3;

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


-- name: GetPlayersByMatchId :many
SELECT 	
	p.id, p.accountid, p.score, p.isready, p.art,
	array(	
		(select m.* from 
		matchcard m where
		m.currowner=p.id and m.state='Hand'
		)) as hand,
	array(	
		(select m.* from 
		matchcard m where
		m.currowner=p.id and m.state='Play'
		)) as board,
	array(	
		(select m.* from 
		matchcard m where
		m.currowner=p.id and m.state='Kitty'
		)) as kitty
FROM player as p
LEFT JOIN
    match_player ON p.id=match_player.playerid
WHERE
    match_player.matchid = $1;



-- name: GetPlayerJSON :many
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
            WHERE m.currowner = p.id AND m.state = 'Play'
        )
    )
FROM player as p
LEFT JOIN
    match_player ON p.id=match_player.playerid
WHERE
    match_player.matchid = $1;

