package store

import (
	"github.com/bardic/gocrib/queries/queries"
	"github.com/labstack/echo/v4"
)

type CardStore struct {
	Store
}

func NewCardStore(q *queries.Queries, c echo.Context) *CardStore {
	return &CardStore{
		Store: Store{
			q: q,
			c: c,
		},
	}
}

func (p *CardStore) GetCardsForPlayerIdFromDeckId(params queries.GetCardsForPlayerIdFromDeckIdParams) ([]queries.GetCardsForPlayerIdFromDeckIdRow, error) {
	cards, err := p.q.GetCardsForPlayerIdFromDeckId(p.c.Request().Context(), params)

	if err != nil {
		return nil, err
	}

	return cards, nil
}

func (p *CardStore) GetCardsForMatchIdAndState(params queries.GetCardsForMatchIdAndStateParams) ([]queries.GetCardsForMatchIdAndStateRow, error) {

	cards, err := p.q.GetCardsForMatchIdAndState(p.c.Request().Context(), params)

	if err != nil {
		return nil, err
	}

	return cards, nil
}

func (p *CardStore) UpdateMatchCardState(params queries.UpdateMatchCardStateParams) error {
	return p.q.UpdateMatchCardState(p.c.Request().Context(), params)
}
