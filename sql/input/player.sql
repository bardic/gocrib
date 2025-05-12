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

-- name: PlayerJoinMatch :exec
INSERT INTO 
    match_player (matchid, playerid)     
VALUES 
    ($1, $2);

-- name: UpdatePlayer :one
UPDATE player SET 
		score = $1, 
		isReady = $2,
		art = $3 
	WHERE 
		id = $4
    RETURNING *;

-- name: UpdatePlayerReady :exec
UPDATE player SET isReady = $1 WHERE id = $2;

-- name: UpdatePlayerTurnOrder :exec
UPDATE match_player SET turnorder = $1 WHERE matchid = $2 AND playerid = $3;

-- name: GetPlayerById :one
SELECT player.* FROM player WHERE id=$1 LIMIT 1;

-- name: GetPlayerByMatchAndAccountId :one
SELECT 
    player.*
FROM
    player
LEFT JOIN
    match_player ON player.id=match_player.playerid
WHERE
    match_player.matchid = $1 AND player.accountid = $2;