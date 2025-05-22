package store

import (
	"github.com/bardic/gocrib/vo"
	"github.com/labstack/echo/v4"
)

type DeckStore struct {
	Store
}

func (p *DeckStore) CreateDeck(ctx echo.Context, matchID int) (*vo.Deck, error) {
	deck, err := p.q().CreateDeck(ctx.Request().Context(), matchID)

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

func (p *DeckStore) ResetDeckForMatchID(ctx echo.Context, matchID int) error {
	err := p.q().ResetDeckForMatchId(ctx.Request().Context(), matchID)

	defer p.Close()

	if err != nil {
		return err
	}

	return nil
}
