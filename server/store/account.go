package store

import (
	"github.com/bardic/gocrib/vo"
	"github.com/labstack/echo/v4"
)

type AccountStore struct {
	Store
}

func (a *AccountStore) GetAccountByID(ctx echo.Context, id int) (*vo.Account, error) {
	account, err := a.q().GetAccountForId(ctx.Request().Context(), id)
	defer a.Close()

	if err != nil {
		return nil, err
	}

	return &vo.Account{
		ID:   account.ID,
		Name: account.Name,
	}, nil
}
