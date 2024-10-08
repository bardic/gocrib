package match

import (
	"context"
	"net/http"

	conn "github.com/bardic/cribbage/server/db"
	"github.com/bardic/cribbage/server/utils"
	"github.com/bardic/cribbagev2/model"
	"github.com/labstack/echo/v4"
)

// Create godoc
// @Summary      Create new match
// @Description
// @Tags         match
// @Accept       json
// @Produce      json
// @Param details body model.MatchRequirements true "MatchRequirements"
// @Success      200  {object}  model.GameMatch
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
		return err
	}

	match := model.GameMatch{}
	match.DeckId = d.Id

	match.PlayerIds = []int{details.RequesterId}
	match.EloRangeMin = details.EloRangeMin
	match.EloRangeMax = details.EloRangeMax
	match.PrivateMatch = details.IsPrivate
	match.GameState = model.WaitingState

	args := utils.ParseMatch(match)

	query := `INSERT INTO match(
				playerIds,
				privateMatch,
				eloRangeMin,
				eloRangeMax,
				deckId,
				cutGameCardId,
				currentPlayerTurn,
				turnPassTimestamps,
				gameState,
				art)
			VALUES (
				@playerIds,
				@privateMatch,
				@eloRangeMin,
				@eloRangeMax,
				@deckId,
				@cutGameCardId,
				@currentPlayerTurn,
				@turnPassTimestamps,
				@gameState,
				@art)
			RETURNING id`

	var matchId int
	err = db.QueryRow(
		context.Background(),
		query,
		args).Scan(&matchId)

	if err != nil {
		return err
	}

	// d = *d.Shuffle()
	// for i := 0; i < 12; i++ {
	// 	if i%2 == 0 {
	// 		match.Players[0].Hand = append(match.Players[0].Hand, d.Cards[i].CardId)
	// 	} else {
	// 		match.Players[1].Hand = append(match.Players[1].Hand, d.Cards[i].CardId)
	// 	}
	// }

	// match.GameState = model.DiscardState

	// // err = updateMatch(match)

	// if err != nil {
	// 	return err
	// }

	// if match.Players[0].Id == details.RequesterId {
	// 	match.Players[1].Hand = []int{}
	// } else {
	// 	match.Players[0].Hand = []int{}
	// }

	match.Id = matchId
	return c.JSON(http.StatusOK, match)
}
