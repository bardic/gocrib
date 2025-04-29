package store

import (
	"github.com/bardic/gocrib/queries/queries"
	"github.com/labstack/echo/v4"
)

type CardStore struct {
	Store
}

func (p *CardStore) GetCardsForPlayerIdFromDeckId(ctx echo.Context, params queries.GetCardsForPlayerIdFromDeckIdParams) ([]queries.GetCardsForPlayerIdFromDeckIdRow, error) {
	cards, err := p.q().GetCardsForPlayerIdFromDeckId(ctx.Request().Context(), params)
	defer p.Close()

	if err != nil {
		return nil, err
	}

	return cards, nil
}

func (p *CardStore) GetCardsForMatchIdAndState(ctx echo.Context, params queries.GetCardsForMatchIdAndStateParams) ([]queries.GetCardsForMatchIdAndStateRow, error) {
	cards, err := p.q().GetCardsForMatchIdAndState(ctx.Request().Context(), params)
	defer p.Close()

	if err != nil {
		return nil, err
	}

	return cards, nil
}

func (p *CardStore) UpdateMatchCardState(ctx echo.Context, params queries.UpdateMatchCardStateParams) error {
	err := p.q().UpdateMatchCardState(ctx.Request().Context(), params)
	defer p.Close()
	if err != nil {
		return err
	}

	return nil
}

func (p *CardStore) GetCards(ctx echo.Context) ([]queries.Card, error) {
	cards, err := p.q().GetCards(ctx.Request().Context())
	defer p.Close()

	if err != nil {
		return nil, err
	}

	return cards, nil
}

func (p *CardStore) CreateMatchCard(ctx echo.Context, params queries.CreateMatchCardParams) (queries.Matchcard, error) {
	mc, err := p.q().CreateMatchCard(ctx.Request().Context(), queries.CreateMatchCardParams{
		Cardid:    params.Cardid,
		Origowner: params.Origowner,
		Currowner: params.Currowner,
		State:     params.State,
	})

	defer p.Close()

	if err != nil {
		return queries.Matchcard{}, err
	}

	return mc, nil
}
