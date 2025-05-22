package route

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

// Login route
//
//	@Summary				Login in with an accountId
//	@Description 		Login route for account - takes an account id and returns the account details
//	@Tags						account
//	@Accept					json
//	@Produce				json
//	@Param					accountId		path			int	true	"account id"'
//	@Success				200					{object}	vo.Account
//	@Failure				400					{object}	error
//	@Failure				500					{object}	error
//	@Router					/account/{accountId} 	[post]
func (h *Handler) Login(c echo.Context) error {
	accountID, err := strconv.Atoi(c.Param("accountId"))
	if err != nil {
		return h.BadParams(c, "can't convert param accountId", err)
	}

	account, err := h.AccountStore.GetAccountByID(c, accountID)
	if err != nil {
		return h.InternalError(c, "failed to get account by id", err)
	}

	return h.Ok(c, account)
}
