package player

import (
	"net/http"
	"strconv"

	"github.com/bardic/gocrib/queries/queries"
	conn "github.com/bardic/gocrib/server/db"
	"github.com/bardic/gocrib/vo"
	"github.com/labstack/echo/v4"
)

// Create godoc
//
//	@Summary	Get player by barcode
//	@Description
//	@Tags		players
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int	true	"search for match by id"'
//	@Param		matchId	path		int	true	"search for match by id"'
//	@Success	200	{object}	vo.GamePlayer
//	@Failure	400	{object}	error
//	@Failure	404	{object}	error
//	@Failure	500	{object}	error
//	@Router		/match/{matchId}/player/{id} [get]
func GetPlayer(c echo.Context) error {
	id := c.Param("id")
	//match := c.Param("matchId")

	p1Id, err := strconv.Atoi(id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)
	ctx := c.Request().Context()

	p, err := q.GetPlayer(ctx, &p1Id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	hand, err := q.GetMarchCardsByType(ctx, queries.GetMarchCardsByTypeParams{
		ID:    &p1Id,
		State: queries.CardstateHand,
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	play, err := q.GetMarchCardsByType(ctx, queries.GetMarchCardsByTypeParams{
		ID:    &p1Id,
		State: queries.CardstatePlay,
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	kitty, err := q.GetMarchCardsByType(ctx, queries.GetMarchCardsByTypeParams{
		ID:    &p1Id,
		State: queries.CardstateKitty,
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	gamePlayer := vo.GamePlayer{
		Player: p,
		Hand:   hand,
		Play:   play,
		Kitty:  kitty,
	}

	return c.JSON(http.StatusOK, gamePlayer)
}
