package player

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/bardic/gocrib/queries/queries"
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
//	@Param		matchId		path		int	true	"match id"'
//	@Param		playerId	path		int	true	"from player id"'
//	@Param		toPlayerId	path		int	true	"to player id"'
// 	@Param		details	body		vo.HandModifier	true	"HandModifier object"
//	@Success	200		{object}	queries.Match
//	@Failure	400		{object}	error
//	@Failure	404		{object}	error
//	@Failure	500		{object}	error
//	@Router		/match/{matchId}/player/{fromPlayerId}/to/{toPlayerId}/play [put]

func (h *PlayerHandler) UpdatePlay(c echo.Context) error {
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
			State:     queries.CardstatePlay,
			Origowner: &fromPlayerID,
			Currowner: &toPlayerID,
		})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
	}

	match, err := h.MatchStore.GetMatch(c, &matchID)

	for i, player := range match.Players {
		if *player.ID == toPlayerID {
			playerIndex := i
			nextPlayer := match.Players[(playerIndex+1)%len(match.Players)]

			err = h.MatchStore.UpateMatchCurrentPlayerTurn(c, queries.UpateMatchCurrentPlayerTurnParams{
				ID:                &matchID,
				Currentplayerturn: nextPlayer.ID,
			})
			if err != nil {
				fmt.Println(err)
			}
			break
		}
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, match)
}
