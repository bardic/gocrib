package player

import (
	"net/http"
	"strconv"

	"github.com/bardic/gocrib/queries/queries"
	"github.com/bardic/gocrib/vo"

	"github.com/labstack/echo/v4"
)

// Create godoc
//
//	@Summary	Update kitty with ids
//	@Description
//	@Tags		match
//	@Accept		json
//	@Produce	json
//	@Param		matchId		path		int	true	"match id"'
//	@Param		playerId	path		int	true	"from player id"'
//	@Param		toPlayerId	path		int	true	"to player id"'
//	@Param		details	body		vo.HandModifier	true	"array of ids to add to kitty"
//	@Success	200		{object}	queries.Match
//	@Failure	400		{object}	error
//	@Failure	404		{object}	error
//	@Failure	500		{object}	error
//	@Router		/match/{matchId}/player/{fromPlayerId}/to/{toPlayerId}/kitty [put]
func (h *PlayerHandler) UpdateKitty(c echo.Context) error {
	matchId, err := strconv.Atoi(c.Param("matchId"))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	fromPlayerId, err := strconv.Atoi(c.Param("playerId"))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	toPlayerId, err := strconv.Atoi(c.Param("toPlayerId"))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	details := &vo.HandModifier{}
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	for _, cardId := range details.CardIds {
		err := h.cardStore.UpdateMatchCardState(queries.UpdateMatchCardStateParams{
			ID:        &cardId,
			State:     queries.CardstateKitty,
			Origowner: &fromPlayerId,
			Currowner: &toPlayerId,
		})

		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
	}

	err = h.playerStore.UpdatePlayerReady(queries.UpdatePlayerReadyParams{
		Isready: true,
		ID:      &fromPlayerId,
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	match, err := h.matchStore.GetMatch(&matchId)

	allReady := true
	for _, player := range match.Players {
		if !player.Isready {
			allReady = false
			break
		}
	}

	if allReady {
		_, err = h.matchStore.UpdateMatchState(queries.UpdateMatchStateParams{
			Gamestate: queries.GamestateCut,
			ID:        &matchId,
		})

		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
	}

	return c.JSON(http.StatusOK, match)
}
