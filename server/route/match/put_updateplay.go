package match

import (
	"context"
	"net/http"

	"queries"
	conn "server/db"
	"server/utils"

	"vo"

	"github.com/labstack/echo/v4"
)

// Create godoc
// @Summary      Update play with ids
// @Description
// @Tags         match
// @Accept       json
// @Produce      json
// @Param details body vo.HandModifier true "array of ids to add to play"
// @Success      200  {object}  queries.Match
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /player/play [put]
func UpdatePlay(c echo.Context) error {
	details := &vo.HandModifier{}
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	m, err := utils.UpdatePlay(*details)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	err = utils.UpdateGameState(m.ID, queries.GamestateOpponentState)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx := context.Background()

	err = q.PassTurn(ctx, m.ID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, m)
}
