package route

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// JoinMatch route
//
//	@Summary			Join match by id
//	@Description
//	@Tags					match
//	@Accept				json
//	@Produce			json
//	@Param				matchId			path			int				true	"match id"'
//	@Param				accountId		path			int				true	"account id"'
//	@Success			200					{object}	vo.Match
//	@Failure			400					{object}	error
//	@Failure			500					{object}	error
//	@Router				/match/{matchId}/join/{accountId} [put]
func (h *Handler) JoinMatch(c echo.Context) error {
	matchID, err := strconv.Atoi(c.Param("matchId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	accountID, err := strconv.Atoi(c.Param("accountId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	player, err := h.PlayerStore.CreatePlayer(c, accountID)
	if err != nil {
		return h.InternalError(c, "", err)
	}

	err = h.PlayerStore.PlayerJoinMatch(c, matchID, player.ID)
	if err != nil {
		return h.InternalError(c, "", err)
	}

	_, err = h.MatchStore.GetMatch(c, matchID)
	if err != nil {
		return h.InternalError(c, "", err)
	}

	match, err := h.MatchStore.GetMatch(c, matchID)
	if err != nil {
		return err
	}

	turnOrder := 1
	for _, player := range match.Players {
		err = h.MatchStore.UpdatePlayerTurnOrder(c, turnOrder, matchID, player.ID)
		if err != nil {
			return err
		}

		if turnOrder == 1 {
			err = h.MatchStore.UpateMatchCurrentPlayerTurn(c, matchID, player.ID)
			if err != nil {
				return err
			}

			err = h.MatchStore.UpdateDealerForMatch(c, matchID, player.ID)
			if err != nil {
				return err
			}
		}

		turnOrder++
	}

	// numOfPlayers := len(match.Players)

	// cardsToDealPerPlayer := 6
	// if numOfPlayers > 2 {
	// 	cardsToDealPerPlayer = 5
	// }

	// Deal cards

	// cards, err := h.MatchStore.GetCardsForMatchID(c, matchId)
	// if err != nil {
	// 	return err
	// }
	//
	// for i := range numOfPlayers {
	// 	for j := range cardsToDealPerPlayer {
	// 		card := cards[(i*cardsToDealPerPlayer)+j]
	// 		h.CardStore.UpdateMatchCardState(c, queries.UpdateMatchCardStateParams{
	// 			State:     queries.CardstateHand,
	// 			Origowner: match.Players[i].ID,
	// 			Currowner: match.Players[i].ID,
	// 			ID:        card.ID_2,
	// 		})
	//
	// 	}
	// }

	_, err = h.MatchStore.UpdateMatchState(c, matchID, "Discard")
	if err != nil {
		return h.InternalError(c, "", err)
	}

	return c.JSON(http.StatusOK, match)
}
