package player

import (
	"context"
	"net/http"
	"time"

	"queries"

	conn "github.com/bardic/gocrib/server/db"

	"github.com/labstack/echo/v4"
)

// Create godoc
//
//	@Summary	Create new player
//	@Description
//	@Tags		players
//	@Accept		json
//	@Produce	json
//	@Param		details	body		int	true	"player Object to save"
//	@Success	200		{object}	queries.Player
//	@Failure	400		{object}	error
//	@Failure	404		{object}	error
//	@Failure	500		{object}	error
//	@Router		/match/player/ [post]
func NewPlayer(c echo.Context) error {
	id := new(int32)
	if err := c.Bind(id); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	p, err := NewPlayerQuery(*id)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, p)
}

func NewPlayerQuery(accountId int32) (*queries.Player, error) {

	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	p, err := q.CreatePlayer(ctx, queries.CreatePlayerParams{
		Accountid: accountId,
		Hand:      []int32{},
		Kitty:     []int32{},
		Play:      []int32{},
		Score:     0,
		Isready:   false,
		Art:       "default.png",
	})

	if err != nil {
		return nil, err
	}

	return &p, nil
}
