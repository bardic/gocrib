package store

import (
	"github.com/bardic/gocrib/queries/queries"
	"github.com/labstack/echo/v4"
)

type DeckStore struct {
	Store
}

func NewDeckStore(q *queries.Queries, c echo.Context) *DeckStore {
	return &DeckStore{
		Store: Store{
			q: q,
			c: c,
		},
	}
}

func (p *DeckStore) GetDeckForMatchId(matchId *int) (*queries.Deck, error) {
	deck, err := p.q.GetDeckForMatchId(p.c.Request().Context(), matchId)

	if err != nil {
		return nil, err
	}

	return &deck, nil
}
