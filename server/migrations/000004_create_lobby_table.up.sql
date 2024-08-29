CREATE TABLE IF NOT EXISTS lobby (
  id int8 PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
  accountIds integer [] NOT NULL,
  creatationDate timestamptz NOT NULL DEFAULT now(),
  privateMatch boolean NOT NULL,
  eloRangeMin int NOT NULL,
  eloRangeMax int NOT NULL
);