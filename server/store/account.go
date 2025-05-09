package store

import (
	"github.com/bardic/gocrib/queries/queries"
	"github.com/labstack/echo/v4"
)

type AccountStore struct {
	Store
}

func (a *AccountStore) GetAccountByID(ctx echo.Context, id *int) (*queries.Account, error) {
	account, err := a.q().GetAccountForId(ctx.Request().Context(), id)
	defer a.Close()

	if err != nil {
		return nil, err
	}

	return &account, nil
}
