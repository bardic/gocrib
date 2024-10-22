package match

import (
	"net/http"

	"github.com/bardic/gocrib/model"
	"github.com/bardic/gocrib/queries"
	"github.com/bardic/gocrib/server/utils"
	"github.com/labstack/echo/v4"
)

// Create godoc
// @Summary      Join match by id
// @Description
// @Tags         match
// @Accept       json
// @Produce      json
// @Param details body model.JoinMatchReq true "match Object to update"
// @Success      200  {object}  model.MatchDetailsResponse
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /player/match/join [put]
func JoinMatch(c echo.Context) error {
	details := new(model.JoinMatchReq)
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	p, err := utils.GetPlayerById(details.PlayerId)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	details.PlayerId = int(p.ID)

	err = utils.UpdatePlayersInMatch(*details)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, model.MatchDetailsResponse{
		MatchId:   int(details.MatchId),
		PlayerId:  int(p.ID),
		GameState: queries.GamestateJoinGameState,
	})
}
