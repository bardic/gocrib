-- name: CreateDeck :one
INSERT INTO deck(cutmatchcardid) VALUES (null) RETURNING *;

-- name: AddCardMatchToDeck :exec
INSERT INTO deck_matchcard (deckid, matchcardid) VALUES ($1, $2);

-- name: GetDeckForMatchId :one
SELECT deck.* FROM deck
LEFT JOIN
    match ON deck.id=match.deckid
WHERE match.id=$1 LIMIT 1;

-- name: GetCardsForDeckId :many
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