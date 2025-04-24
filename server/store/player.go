package store

import (
	"github.com/bardic/gocrib/queries/queries"
	"github.com/labstack/echo/v4"
)

type PlayerStore struct {
	Store
}

func NewPlayerStore(q *queries.Queries, c echo.Context) *PlayerStore {
	return &PlayerStore{
		Store: Store{
			q: q,
			c: c,
		},
	}
}

// func (p *PlayerStore) CreatePlayer(params queries.CreatePlayerParams) (*queries.Player, error) {
// 	player, err := p.q.GetPlayerById(p.c.Request().Context(), params.Accountid)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return &player, nil
// }

func (p *PlayerStore) GetPlayerById(id *int) (*queries.Player, error) {
	player, err := p.q.GetPlayerById(p.c.Request().Context(), id)

	if err != nil {
		return nil, err
	}

	return &player, nil
}

func (p *PlayerStore) UpdatePlayerReady(params queries.UpdatePlayerReadyParams) error {
	return p.q.UpdatePlayerReady(p.c.Request().Context(), params)
}
