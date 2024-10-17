CREATE TABLE IF NOT EXISTS match (
  id integer PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
  playerIds integer [] NOT NULL,
  creationDate timestamptz NOT NULL DEFAULT now(),
  privateMatch boolean NOT NULL,
  eloRangeMin int NOT NULL,
  eloRangeMax int NOT NULL,
  deckId integer NOT NULL,
  cutGameCardId integer NOT NULL,
  currentPlayerTurn integer NOT NULL,
  turnPassTimestamps timestamptz [] NOT NULL,
  gameState integer NOT NULL,
  art varchar(256) NOT NULL
);