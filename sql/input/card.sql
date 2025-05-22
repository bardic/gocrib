-- name: CreateMatchCard :copyfrom
INSERT INTO matchcard (cardid, origowner, currowner, state, deckid) VALUES ($1, $2, $3, $4, $5);

-- name: UpdateMatchCardState :exec
UPDATE matchcard SET state = $1, origowner = $2, currowner = $3 WHERE id = $4;

-- name: GetCards :many
SELECT * FROM card;

-- name: GetCardsForPlayerIdFromDeckId :many
SELECT
    card.*,
    matchcard.* 
FROM
    card
LEFT JOIN
    matchcard ON card.id=matchcard.cardid
WHERE
    matchcard.deckid = $1 AND matchcard.origowner = $2;

-- name: GetCardsForMatchIdAndState :many
SELECT 
    deck.*,
    matchcard.*,
    card.*
FROM
    matchcard
LEFT JOIN
    deck ON matchcard.deckid=deck.id
LEFT JOIN
    match ON deck.matchId=match.id
left join 
	card on card.id=matchcard.cardid
WHERE
    match.id = $1 AND matchcard.state = $2;

-- name: GetCardsForMatchId :many
SELECT 
    deck.*,
    matchcard.*,
    card.*
FROM
    matchcard
LEFT JOIN
    deck ON matchcard.deckid=deck.id
LEFT JOIN
    match ON deck.matchId=match.id
left join 
	card on card.id=matchcard.cardid
WHERE
    match.id = $1;
