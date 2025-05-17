-- name: CreateDeck :one
INSERT INTO deck(cutmatchcardid) VALUES (null) RETURNING *;

-- name: AddCardToDeck :copyfrom
INSERT INTO deck_matchcard (deckid, matchcardid) VALUES ($1, $2);


