-- name: CreateDeck :one
INSERT INTO deck(cutmatchcardid) VALUES (null) RETURNING *;

-- name: AddCardToDeck :copyfrom
INSERT INTO deck_matchcard (deckid, matchcardid) VALUES ($1, $2);

-- name: GetDeckForMatchId :one
SELECT deck.* FROM deck
LEFT JOIN
    match ON deck.id=match.deckid
WHERE match.id=$1 LIMIT 1;
