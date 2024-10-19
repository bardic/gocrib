package card

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Create godoc
// @Summary      Get card by barcode
// @Description
// @Tags         cards
// @Accept       json
// @Produce      json
// @Param        barcode    query     string  true  "search for card by barcode"'
// @Param        storeId    query     string  true  "Store in which the barcode was found"'
// @Success      200  {object}  queries.Card
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /card/ [delete]
func DeleteCard(c echo.Context) error {

	return c.JSON(http.StatusOK, nil)
}
