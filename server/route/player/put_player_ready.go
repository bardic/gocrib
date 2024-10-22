package player

import (
	"context"
	"net/http"

	"github.com/bardic/gocrib/queries"
	conn "github.com/bardic/gocrib/server/db"
	"github.com/bardic/gocrib/server/utils"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

// Create godoc
// @Summary      Update player to mark as ready
// @Description
// @Tags         players
// @Accept       json
// @Produce      json
// @Param details body int true "player id to update"
// @Success      200  {object}  bool
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /player/ready [put]
func PlayerReady(c echo.Context) error {
	details := new(int)
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx := context.Background()

	_, err := ReadyPlayerById(c, *details)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	match, err := utils.GetMatchForPlayerId(*details)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	if len(match.Players) != 2 {
		return c.JSON(http.StatusOK, nil)
	}

	deck, err := utils.NewDeck()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	for _, id := range deck.Cards {
		_, err := q.CreateMatchCards(ctx, queries.CreateMatchCardsParams{
			MatchID:   match.ID,
			Cardid:    id,
			State:     queries.CardstateDeck,
			Origowner: pgtype.Int4{Int32: 0},
			Currowner: pgtype.Int4{Int32: 0},
		})

		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
	}

	err = q.UpdateMatchWithDeckId(ctx, queries.UpdateMatchWithDeckIdParams{
		ID:     match.ID,
		Deckid: int32(deck.ID),
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	match.Deckid = deck.ID

	if utils.PlayersReady(match.Players) {
		utils.Deal(match)
		utils.UpdateGameState(int(match.ID), queries.GamestateDiscardState)
	}

	return c.JSON(http.StatusOK, nil)
}

func ReadyPlayerById(c echo.Context, playerId int) (bool, error) {
	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx := context.Background()

	err := q.UpdatePlayerReady(ctx, queries.UpdatePlayerReadyParams{
		ID:      int32(playerId),
		Isready: true,
	})

	if err != nil {
		return false, err
	}

	return true, nil
}
