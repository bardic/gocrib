package deck

import (
	"net/http"
	"strconv"

	"github.com/bardic/gocrib/server/utils"
	"github.com/labstack/echo/v4"
)

// Create godoc
// @Summary      Get deck by id
// @Description
// @Tags         deck
// @Accept       json
// @Produce      json
// @Param        id    query     string  true  "search for deck by id"'
// @Success      200  {object}  model.GameDeck
// @Failure      404  {object}  error
// @Failure      422  {object}  error
// @Router       /player/match/deck/ [get]
func GetDeck(c echo.Context) error {
	p := c.Request().URL.Query().Get("id")
	id, err := strconv.Atoi(p)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	deck, err := utils.GetDeckById(id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, deck)
}
