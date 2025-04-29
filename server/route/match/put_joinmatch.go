package match

import (
	"net/http"
	"strconv"

	"github.com/bardic/gocrib/queries/queries"

	"github.com/labstack/echo/v4"
)

// Return match details if the player is able to join the match
//
//	@Summary	Join match by id
//	@Description
//	@Tags		match
//	@Accept		json
//	@Produce	json
//	@Param		matchId	path		int	true	"match id"'
//	@Param		accountId	path		int	true	"account id"'
//	@Success	200		{object}	vo.MatchDetailsResponse
//	@Failure	400		{object}	error
//	@Failure	404		{object}	error
//	@Failure	500		{object}	error
//	@Router		/match/{matchId}/join/{accountId} [put]
func (h *MatchHandler) JoinMatch(c echo.Context) error {
	matchId, err := strconv.Atoi(c.Param("matchId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	accountId, err := strconv.Atoi(c.Param("accountId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	score := 0

	player, err := h.PlayerStore.CreatePlayer(c, queries.CreatePlayerParams{
		Accountid: &accountId,
		Score:     &score,
		Isready:   false,
		Art:       "default.png",
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	err = h.BaseHandler.PlayerStore.PlayerJoinMatch(c, queries.PlayerJoinMatchParams{
		Matchid:  &matchId,
		Playerid: player.ID,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	_, err = GetMatch(&matchId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	match, err := h.MatchStore.GetMatch(c, &matchId)
	if err != nil {
		return err
	}

	turnOrder := 1
	for _, player := range match.Players {
		h.MatchStore.UpdatePlayerTurnOrder(c, queries.UpdatePlayerTurnOrderParams{
			Turnorder: &turnOrder,
			Matchid:   &matchId,
			Playerid:  player.ID,
		})

		if turnOrder == 1 {
			err = h.MatchStore.UpateMatchCurrentPlayerTurn(c, queries.UpateMatchCurrentPlayerTurnParams{
				ID:                &matchId,
				Currentplayerturn: player.ID,
			})
			if err != nil {
				return err
			}

			err = h.MatchStore.UpdateDealerForMatch(c, queries.UpdateDealerForMatchParams{
				ID:       &matchId,
				Dealerid: player.ID,
			})
			if err != nil {
				return err
			}
		}

		turnOrder++
	}

	numOfPlayers := len(match.Players)

	cardsToDealPerPlayer := 6
	if numOfPlayers > 2 {
		cardsToDealPerPlayer = 5
	}

	// Deal cards

	cards, err := h.MatchStore.GetCardsForMatchId(c, matchId)
	if err != nil {
		return err
	}

	for i := range numOfPlayers {
		for j := range cardsToDealPerPlayer {
			card := cards[(i*cardsToDealPerPlayer)+j]
			h.CardStore.UpdateMatchCardState(c, queries.UpdateMatchCardStateParams{
				State:     queries.CardstateHand,
				Origowner: match.Players[i].ID,
				Currowner: match.Players[i].ID,
				ID:        card.ID_2,
			})

		}
	}

	_, err = h.MatchStore.UpdateMatchState(c, queries.UpdateMatchStateParams{
		ID:        &matchId,
		Gamestate: queries.GamestateDiscard,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, match)
}
