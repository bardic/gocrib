package deck

import (
	"net/http"
	"strconv"

	"github.com/bardic/gocrib/server/controller"
	"github.com/labstack/echo/v4"
)

// Returns the deck for a match id
//
//	@Summary	Get deck by match id
//	@Description
//	@Tags		deck
//	@Accept		json
//	@Produce	json
//	@Param		id	query		string	true	"search for deck by match id"'
//	@Success	200	{object}	vo.GameDeck
//	@Failure	404	{object}	error
//	@Failure	422	{object}	error
//	@Router		/player/match/deck/ [get]
func GetDeckByMatchId(c echo.Context) error {
	p := c.Request().URL.Query().Get("id")
	id, err := strconv.Atoi(p)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	deck, err := controller.GetDeckByMatchId(int32(id))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, deck)
}
