package store

import (
	"encoding/json"

	"github.com/bardic/gocrib/queries/queries"
	"github.com/bardic/gocrib/vo"
	"github.com/labstack/echo/v4"
)

type MatchStore struct {
	Store
}

func NewMatchStore(q *queries.Queries, c echo.Context) *MatchStore {
	return &MatchStore{
		Store: Store{},
	}
}

func (p *MatchStore) GetMatch(ctx echo.Context, matchId *int) (*vo.GameMatch, error) {
	m, err := p.q().GetMatchById(ctx.Request().Context(), matchId)

	defer p.Close()

	if err != nil {
		return nil, err
	}

	var match *vo.GameMatch
	err = json.Unmarshal(m, &match)
	if err != nil {
		return nil, err
	}

	return match, nil
}

func (p *MatchStore) UpdateMatchState(ctx echo.Context, params queries.UpdateMatchStateParams) (queries.Match, error) {
	match, err := p.q().UpdateMatchState(ctx.Request().Context(), params)

	if err != nil {
		return queries.Match{}, err
	}

	defer p.Close()

	return match, nil
}

func (p *MatchStore) UpateMatchCurrentPlayerTurn(ctx echo.Context, params queries.UpateMatchCurrentPlayerTurnParams) error {
	err := p.q().UpateMatchCurrentPlayerTurn(ctx.Request().Context(), params)
	defer p.Close()
	if err != nil {
		return err
	}

	return nil
}

func (p *MatchStore) GetCardsForMatchId(ctx echo.Context, id int) ([]queries.GetCardsForMatchIdRow, error) {
	cards, err := p.q().GetCardsForMatchId(ctx.Request().Context(), &id)

	defer p.Close()

	if err != nil {
		return nil, err
	}

	return cards, nil
}

func (p *MatchStore) CreateMatch(ctx echo.Context, params queries.CreateMatchParams) (queries.Match, error) {
	m, err := p.q().CreateMatch(ctx.Request().Context(), params)

	defer p.Close()

	if err != nil {
		return queries.Match{}, err
	}

	return m, nil
}

func (p *MatchStore) PlayerJoinMatch(ctx echo.Context, params queries.PlayerJoinMatchParams) error {
	err := p.q().PlayerJoinMatch(ctx.Request().Context(), params)

	defer p.Close()

	if err != nil {
		return err
	}

	return nil
}

func (p *MatchStore) UpdateMatchCut(ctx echo.Context, params queries.UpdateMatchCutParams) error {
	err := p.q().UpdateMatchCut(ctx.Request().Context(), params)
	defer p.Close()

	if err != nil {
		return err
	}

	return nil
}

func (p *MatchStore) UpdatePlayerTurnOrder(ctx echo.Context, params queries.UpdatePlayerTurnOrderParams) error {
	err := p.q().UpdatePlayerTurnOrder(ctx.Request().Context(), params)
	defer p.Close()
	if err != nil {
		return err
	}

	return nil
}

func (p *MatchStore) UpdateDealerForMatch(ctx echo.Context, params queries.UpdateDealerForMatchParams) error {
	err := p.q().UpdateDealerForMatch(ctx.Request().Context(), params)
	defer p.Close()
	if err != nil {
		return err
	}
	return nil
}
