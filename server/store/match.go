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
		Store: Store{
			q: q,
			c: c,
		},
	}
}

func (p *MatchStore) GetMatch(matchId *int) (*vo.GameMatch, error) {
	m, err := p.q.GetMatchById(p.c.Request().Context(), matchId)

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

func (p *MatchStore) UpdateMatchState(params queries.UpdateMatchStateParams) (queries.Match, error) {
	return p.q.UpdateMatchState(p.c.Request().Context(), params)
}
