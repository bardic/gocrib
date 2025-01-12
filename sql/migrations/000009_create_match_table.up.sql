CREATE TYPE GameState AS ENUM (
  'New',
  'Waiting',
  'Ready',
  'Determine',
  'Deal',
  'Discard',
  'Cut',
  'PlayOwn',
  'PlayOpponent',
  'PassTurn',
  'Count',
  'Kitty',
  'Won',
  'Lost'
);


CREATE TABLE IF NOT EXISTS match (
  id integer PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
  creationDate timestamptz NOT NULL DEFAULT now(),
  privateMatch boolean NOT NULL,
  eloRangeMin int NOT NULL,
  eloRangeMax int NOT NULL,
  deckId integer NOT NULL,
  CONSTRAINT fk_deck
    FOREIGN KEY(deckId) 
      REFERENCES deck(id),
  cutGameCardId integer NOT NULL,
  -- CONSTRAINT fk_matchCards
  --   FOREIGN KEY(cutGameCardId) 
  --     REFERENCES matchCards(id),
  currentPlayerTurn integer,
  CONSTRAINT fk_player
    FOREIGN KEY(currentPlayerTurn) 
      REFERENCES player(id),
  turnPassTimestamps timestamptz [] NOT NULL,
  gameState GameState NOT NULL,
  art varchar(256) NOT NULL
);