package store

import (
	"github.com/bardic/gocrib/queries/queries"
	"github.com/labstack/echo/v4"
)

type DeckStore struct {
	Store
}

func NewDeckStore() *DeckStore {
	return &DeckStore{
		Store: Store{},
	}
}

func (p *DeckStore) GetDeckForMatchId(ctx echo.Context, matchId *int) (*queries.Deck, error) {
	deck, err := p.q().GetDeckForMatchId(ctx.Request().Context(), matchId)

	defer p.Close()

	if err != nil {
		return nil, err
	}

	return &deck, nil
}

func (p *DeckStore) CreateDeck(ctx echo.Context) (queries.Deck, error) {
	deck, err := p.q().CreateDeck(ctx.Request().Context())

	defer p.Close()

	if err != nil {
		return queries.Deck{}, err
	}

	return deck, nil
}

func (p *DeckStore) AddCardToDeck(ctx echo.Context, params queries.AddCardToDeckParams) error {
	err := p.q().AddCardToDeck(ctx.Request().Context(), params)

	defer p.Close()

	if err != nil {
		return err
	}

	return nil
}

func (p *DeckStore) ResetDeckForMatchId(ctx echo.Context, matchId *int) error {
	err := p.q().ResetDeckForMatchId(ctx.Request().Context(), matchId)

	defer p.Close()

	if err != nil {
		return err
	}

	return nil
}
