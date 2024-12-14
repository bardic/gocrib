package match

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// DeleteMatch takes an match id and deletes it
//
//	@Summary	Delete a match by id
//	@Description
//	@Tags		match
//	@Accept		json
//	@Produce	json
//	@Param		id	query		string	true	"search for match by barcode"'
//	@Success	200	{object}	queries.Match
//	@Failure	400	{object}	error
//	@Failure	404	{object}	error
//	@Failure	500	{object}	error
//	@Router		/admin/match/ [delete]
func DeleteMatch(c echo.Context) error {
	// b := c.Request().URL.Query().Get("barcode")
	// s := c.Request().URL.Query().Get("storeId")

	// match, _ := getmatch(b, s)

	return c.JSON(http.StatusOK, nil)
}
