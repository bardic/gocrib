CREATE TABLE IF NOT EXISTS deck_matchCard (
  deckId integer NOT NULL,
  matchCardId integer NOT NULL,
  PRIMARY KEY (deckId, matchCardId),
  CONSTRAINT fk_deck
    FOREIGN KEY(deckId) 
      REFERENCES deck(id),
  CONSTRAINT fk_matchCard
    FOREIGN KEY(matchCardId) 
      REFERENCES matchCard(id)
);
