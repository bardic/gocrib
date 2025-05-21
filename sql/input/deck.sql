-- name: CreateDeck :one
INSERT INTO deck(cutmatchcardid) VALUES (null) RETURNING *;

