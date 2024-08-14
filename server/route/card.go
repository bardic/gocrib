package route

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/bardic/cribbage/server/model"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

// Create godoc
// @Summary      Create new card
// @Description
// @Tags         cards
// @Accept       json
// @Produce      json
// @Param details body model.card true "card Object to save"
// @Success      200  {object}  model.Card
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /card/ [post]
func NewCard(c echo.Context) error {
	details := new(model.Card)
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	args := parseCard(*details)
	query := "INSERT INTO cards(value, suit, currentOwner, originalOwner, state, art) VALUES (@value, @suit, @currentOwner, @originalOwner, @state, @art)"

	db := model.Pool()
	defer db.Close()

	_, err := db.Exec(
		context.Background(),
		query,
		args)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "meow")
}

// Create godoc
// @Summary      Update card by barcode
// @Description
// @Tags         cards
// @Accept       json
// @Produce      json
// @Param details body model.card true "card Object to save"
// @Success      200  {object}  model.card
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /card/ [put]
func UpdateCard(c echo.Context) error {
	details := new(model.Card)
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	args := parseCard(*details)

	query := "UPDATE cards SET value=@value, suit=@suit, currentOwner=@currentOwner, originalOwner=@originalOwner, state=@state, art=@art"

	db := model.Pool()
	defer db.Close()

	_, err := db.Exec(
		context.Background(),
		query,
		args)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "meow")
}

// Create godoc
// @Summary      Get card by barcode
// @Description
// @Tags         cards
// @Accept       json
// @Produce      json
// @Param        barcode    query     string  true  "search for card by barcode"'
// @Param        storeId    query     string  true  "Store in which the barcode was found"'
// @Success      200  {object}  model.card
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /card/ [get]
func GetCard(c echo.Context) error {
	id := c.Request().URL.Query().Get("id")

	db := model.Pool()
	defer db.Close()

	rows, err := db.Query(context.Background(), "SELECT * FROM cards WHERE id=$1", id)

	v := []model.Card{}

	for rows.Next() {
		var card model.Card

		err := rows.Scan(&card.Id, &card.Value, &card.Suit, &card.CurrOwner, &card.OrigOwner, &card.State, &card.Art)
		if err != nil {
			return err
		}

		v = append(v, card)
	}

	if err != nil {
		return err
	}

	r, _ := json.Marshal(v)

	return c.JSON(http.StatusOK, string(r))
}

// Create godoc
// @Summary      Get card by barcode
// @Description
// @Tags         cards
// @Accept       json
// @Produce      json
// @Param        barcode    query     string  true  "search for card by barcode"'
// @Param        storeId    query     string  true  "Store in which the barcode was found"'
// @Success      200  {object}  model.card
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /card/ [delete]
func DeleteCard(c echo.Context) error {

	return c.JSON(http.StatusOK, nil)
}

func parseCard(details model.Card) pgx.NamedArgs {
	return pgx.NamedArgs{
		"value":         details.Value,
		"suit":          details.Suit,
		"currentOwner":  details.CurrOwner,
		"originalOwner": details.OrigOwner,
		"state":         details.State,
		"art":           details.Art,
	}
}
