package match

import (
	"net/http"
	"strconv"

	"github.com/bardic/cribbage/server/utils"
	"github.com/bardic/gocrib/model"
	"github.com/labstack/echo/v4"
)

// Create godoc
// @Summary      Cut deck by index
// @Description
// @Tags         match
// @Accept       json
// @Produce      json
// @Param details body model.CutDeckReq true "Deck index that is to become the cut"
// @Success      200  {object}  int
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /player/match/cut [put]
func CutDeck(c echo.Context) error {
	details := new(model.CutDeckReq)
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	m, err := utils.GetMatch(details.MatchId)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	d, err := utils.GetDeckById(m.Id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	cutIndex, err := strconv.Atoi(details.CutIndex)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	card := d.Cards[cutIndex]

	m.CutGameCardId = card.CardId

	err = utils.UpdateCut(m.Id, card.CardId)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	utils.UpdateMatchState(details.MatchId, model.DiscardState)

	return c.JSON(http.StatusOK, details.MatchId)
}
