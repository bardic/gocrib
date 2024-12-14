package match

import (
	"net/http"

	"queries"

	"github.com/bardic/gocrib/server/route/helpers"
	"github.com/labstack/echo/v4"
)

// Updates the current match
//
//	@Summary	Update match details
//	@Description
//	@Tags		match
//	@Accept		json
//	@Produce	json
//	@Param		details	body		queries.Match	true	"match Object to update"
//	@Success	200		{object}	queries.Match
//	@Failure	400		{object}	error
//	@Failure	404		{object}	error
//	@Failure	500		{object}	error
//	@Router		/player/match/ [put]
func UpdateMatch(c echo.Context) error {
	details := new(queries.Match)
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := helpers.UpdateMatch(*details); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, "meow")
}
