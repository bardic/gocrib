package player

import (
	"net/http"
	"strconv"

	"github.com/bardic/gocrib/queries/queries"
	"github.com/bardic/gocrib/vo"

	"github.com/labstack/echo/v4"
)

// UpdateKitty route
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
	matchID, err := strconv.Atoi(c.Param("matchId"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	fromPlayerID, err := strconv.Atoi(c.Param("playerId"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	toPlayerID, err := strconv.Atoi(c.Param("toPlayerId"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	details := &vo.HandModifier{}
	if err = c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	for _, cardID := range details.CardIDs {
		err = h.CardStore.UpdateMatchCardState(c, queries.UpdateMatchCardStateParams{
			ID:        &cardID,
			State:     queries.CardstateKitty,
			Origowner: &fromPlayerID,
			Currowner: &toPlayerID,
		})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
	}

	err = h.PlayerStore.UpdatePlayerReady(c, queries.UpdatePlayerReadyParams{
		Isready: true,
		ID:      &fromPlayerID,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	match, err := h.MatchStore.GetMatch(c, &matchID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	allReady := true
	for _, player := range match.Players {
		if !player.Isready {
			allReady = false
			break
		}
	}

	if allReady {
		_, err = h.MatchStore.UpdateMatchState(c, queries.UpdateMatchStateParams{
			Gamestate: queries.GamestateCut,
			ID:        &matchID,
		})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
	}

	return c.JSON(http.StatusOK, match)
}
