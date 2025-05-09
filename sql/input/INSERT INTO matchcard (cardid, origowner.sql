INSERT INTO matchcard (cardid, origowner, currowner, state) VALUES ($1, $2, $3, $4) RETURNING *;

INSERT INTO 
    deck_matchcard (deckid, matchcardid) 
VALUES 
    ($1, 
    INSERT INTO 
        matchcard (cardid, origowner, currowner, state) 
    VALUES 
        ($1, $2, $3, $4) RETURNING *;
    )
;
