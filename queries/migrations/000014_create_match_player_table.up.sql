CREATE TABLE IF NOT EXISTS match_player (
  matchId integer NOT NULL,
  playerId integer NOT NULL,
  PRIMARY KEY (matchId, playerId),
  CONSTRAINT fk_match
    FOREIGN KEY(matchId) 
      REFERENCES match(id),
  CONSTRAINT fk_player
    FOREIGN KEY(playerId) 
      REFERENCES player(id)
);