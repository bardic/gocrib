package store

import (
	"github.com/bardic/gocrib/queries/queries"
	"github.com/bardic/gocrib/vo"
	"github.com/labstack/echo/v4"
)

type CardStore struct {
	Store
}

func (p *CardStore) GetCardsForMatchIDAndState(
	ctx echo.Context,
	params queries.GetCardsForMatchIdAndStateParams,
) ([]*vo.Card, error) {
	cards, err := p.q().GetCardsForMatchIdAndState(ctx.Request().Context(), params)
	defer p.Close()

	if err != nil {
		return nil, err
	}

	matchCards := make([]*vo.Card, len(cards))
	for i, c := range matchCards {
		matchCards[i] = &vo.Card{
			ID:        c.ID,
			Cardid:    c.Cardid,
			Origowner: c.Origowner,
			Currowner: c.Currowner,
			State:     c.State,
			Name:      c.Name,
			Suit:      c.Suit,
			Art:       c.Art,
		}
	}

	return matchCards, nil
}

func (p *CardStore) UpdateMatchCardState(ctx echo.Context, params queries.UpdateMatchCardStateParams) error {
	err := p.q().UpdateMatchCardState(ctx.Request().Context(), params)
	defer p.Close()
	if err != nil {
		return err
	}

	return nil
}

func (p *CardStore) GetCards(ctx echo.Context) ([]*vo.Card, error) {
	cards, err := p.q().GetCards(ctx.Request().Context())
	defer p.Close()

	if err != nil {
		return nil, err
	}

	matchCards := make([]*vo.Card, len(cards))

	for i, c := range cards {
		matchCards[i] = &vo.Card{
			ID: c.ID,
			// Cardid:    c.Cardid,
			// Origowner: c.Origowner,
			// Currowner: c.Currowner,
			// State:     c.State,
			// Name:      c.Name,
			// Suit:      c.Suit,
			Art: c.Art,
		}
	}

	return matchCards, nil
}

func (p *CardStore) CreateMatchCard(ctx echo.Context, params []queries.CreateMatchCardParams) error {
	_, err := p.q().CreateMatchCard(ctx.Request().Context(), params)

	defer p.Close()

	if err != nil {
		return err
	}

	return nil
}
