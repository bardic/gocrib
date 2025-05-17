package store

import (
	"github.com/bardic/gocrib/queries/queries"
	"github.com/bardic/gocrib/vo"
	"github.com/labstack/echo/v4"
)

type DeckStore struct {
	Store
}

func (p *DeckStore) CreateDeck(ctx echo.Context) (*vo.Deck, error) {
	deck, err := p.q().CreateDeck(ctx.Request().Context())

	defer p.Close()

	if err != nil {
		return nil, err
	}

	return &vo.Deck{
		ID:             deck.ID,
		Cutmatchcardid: deck.Cutmatchcardid,
		Matchid:        deck.Matchid,
		Cards:          []*vo.Card{},
	}, nil
}

func (p *DeckStore) AddCardToDeck(ctx echo.Context, params []queries.AddCardToDeckParams) error {
	_, err := p.q().AddCardToDeck(ctx.Request().Context(), params)

	defer p.Close()

	if err != nil {
		return err
	}

	return nil
}

func (p *DeckStore) ResetDeckForMatchID(ctx echo.Context, matchID int) error {
	err := p.q().ResetDeckForMatchId(ctx.Request().Context(), matchID)

	defer p.Close()

	if err != nil {
		return err
	}

	return nil
}
