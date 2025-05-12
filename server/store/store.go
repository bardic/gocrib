package store

import (
	"github.com/bardic/gocrib/queries/queries"
	conn "github.com/bardic/gocrib/server/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	db *pgxpool.Pool
}

func (s *Store) q() *queries.Queries {
	s.db = conn.Pool()

	q := queries.New(s.db)

	return q
}

func (s *Store) Close() {
	s.db.Close()
}
