package match

import (
	"context"
	"net/http"

	"github.com/bardic/gocrib/model"
	conn "github.com/bardic/gocrib/server/db"
	"github.com/bardic/gocrib/server/route/player"
	"github.com/bardic/gocrib/server/utils"
	"github.com/labstack/echo/v4"
)

// Create godoc
// @Summary      Create new match
// @Description
// @Tags         match
// @Accept       json
// @Produce      json
// @Param details body model.MatchRequirements true "MatchRequirements"
// @Success      200  {object}  int
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /player/match/ [post]
func NewMatch(c echo.Context) error {
	details := new(model.MatchRequirements)
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	db := conn.Pool()
	defer db.Close()

	d, err := utils.NewDeck()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	p, err := player.NewPlayerQuery(details.AccountId)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	match := queries.Match{}
	match.DeckId = d.Id

	match.PlayerIds = []int{p.Id}
	match.EloRangeMin = details.EloRangeMin
	match.EloRangeMax = details.EloRangeMax
	match.PrivateMatch = details.IsPrivate
	match.GameState = model.NewGameState

	args := utils.ParseMatch(match)

	//MEOWCAKES
	// query := `INSERT INTO match(
	// 			playerIds,
	// 			privateMatch,
	// 			eloRangeMin,
	// 			eloRangeMax,
	// 			deckId,
	// 			cutGameCardId,
	// 			currentplayerturn,
	// 			turnPassTimestamps,
	// 			gameState,
	// 			art)
	// 		VALUES (
	// 			@playerIds,
	// 			@privateMatch,
	// 			@eloRangeMin,
	// 			@eloRangeMax,
	// 			@deckId,
	// 			@cutGameCardId,
	// 			@currentPlayerTurn,
	// 			@turnPassTimestamps,
	// 			@gameState,
	// 			@art)
	// 		RETURNING id`

	var matchId int
	err = db.QueryRow(
		context.Background(),
		query,
		args).Scan(&matchId)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, model.MatchDetailsResponse{
		MatchId:   matchId,
		PlayerId:  p.Id,
		GameState: model.NewGameState,
	})
}
