package deck

import (
	"net/http"
	"strconv"

	"github.com/bardic/gocrib/queries/queries"
	"github.com/bardic/gocrib/vo"

	"github.com/labstack/echo/v4"
)

// Create godoc
//
//	@Summary	GetKitty
//	@Description
//	@Tags		match
//	@Accept		json
//	@Produce	json
//	@Param		matchId		path		int	true	"match id"'
//	@Success	200		{object}	vo.Hand
//	@Failure	400		{object}	error
//	@Failure	404		{object}	error
//	@Failure	500		{object}	error
//	@Router		/match/{matchId}/deck/kitty [get]
func (h *DeckHandler) GetKitty(c echo.Context) error {
	matchId, err := strconv.Atoi(c.Param("matchId"))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	cardsResults, err := h.CardStore.GetCardsForMatchIdAndState(c, queries.GetCardsForMatchIdAndStateParams{
		ID:    &matchId,
		State: queries.CardstateKitty,
	})

	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}

	cards := []queries.Matchcard{}
	for _, res := range cardsResults {
		cards = append(cards, queries.Matchcard{
			ID:        res.ID,
			Cardid:    res.Cardid,
			Origowner: res.Origowner,
			Currowner: res.Currowner,
			State:     res.State,
		})
	}

	hand := vo.Hand{
		Cards: cards,
	}

	return c.JSON(http.StatusOK, hand)
}
