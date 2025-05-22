-- name: CreateDeck :one
INSERT INTO deck(cutmatchcardid, matchid) VALUES (-1, $1) RETURNING *;

