package match

import (
	"context"
	"net/http"
	"time"

	"queries"

	"github.com/bardic/gocrib/server/controller"
	conn "github.com/bardic/gocrib/server/db"
	"github.com/bardic/gocrib/server/route/helpers"

	"github.com/bardic/gocrib/vo"

	"github.com/labstack/echo/v4"
)

// Updates a matches cards state
//
//	@Summary	Update play
//	@Description
//	@Tags		match
//	@Accept		json
//	@Produce	json
//	@Param		details	body		vo.HandModifier	true	"HandModifier object"
//	@Success	200		{object}	queries.Match
//	@Failure	400		{object}	error
//	@Failure	404		{object}	error
//	@Failure	500		{object}	error
//	@Router		/play [put]
func UpdatePlay(c echo.Context) error {
	details := &vo.HandModifier{}
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	m, err := controller.UpdatePlay(*details)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	err = helpers.UpdateGameState(m.ID, queries.GamestateOpponentState)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = q.PassTurn(ctx, m.ID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, m)
}
