package store

import (
	"github.com/bardic/gocrib/queries/queries"
	"github.com/bardic/gocrib/vo"
	"github.com/labstack/echo/v4"
)

type PlayerStore struct {
	Store
}

func (p *PlayerStore) CreatePlayer(ctx echo.Context, params queries.CreatePlayerParams) (*vo.Player, error) {
	player, err := p.q().CreatePlayer(ctx.Request().Context(), params)
	defer p.Close()
	if err != nil {
		return nil, err
	}

	return &vo.Player{
		ID:        player.ID,
		Accountid: player.Accountid,
		Score:     player.Score,
		Isready:   player.Isready,
		Art:       player.Art,
	}, nil
}

func (p *PlayerStore) UpdatePlayerReady(ctx echo.Context, params queries.UpdatePlayerReadyParams) error {
	err := p.q().UpdatePlayerReady(ctx.Request().Context(), params)

	defer p.Close()

	if err != nil {
		return err
	}
	return nil
}

func (p *PlayerStore) GetPlayerByMatchAndAccountID(
	ctx echo.Context,
	params queries.GetPlayerByMatchAndAccountIdParams,
) (*vo.Player, error) {
	player, err := p.q().GetPlayerByMatchAndAccountId(ctx.Request().Context(), params)

	defer p.Close()

	if err != nil {
		return nil, err
	}

	return &vo.Player{
		ID:        player.ID,
		Accountid: player.Accountid,
		Score:     player.Score,
		Isready:   player.Isready,
		Art:       player.Art,
	}, nil
}

func (p *PlayerStore) PlayerJoinMatch(ctx echo.Context, params queries.PlayerJoinMatchParams) error {
	err := p.q().PlayerJoinMatch(ctx.Request().Context(), params)
	defer p.Close()
	if err != nil {
		return err
	}
	return nil
}
