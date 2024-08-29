package route

import (
	"context"
	"encoding/json"
	"errors"
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
// @Param details body model.Card true "card Object to save"
// @Success      200  {object}  model.Card
// @Failure      400  {object}  error
// @Failure      500  {object}  error
// @Router       /admin/card/ [post]
func NewCard(c echo.Context) error {
	details := new(model.Card)
	if err := c.Bind(details); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	args := parseCard(*details)
	query := "INSERT INTO cards(value, suit, art) VALUES (@value, @suit, @art)"

	db := model.Pool()
	defer db.Close()

	_, err := db.Exec(
		context.Background(),
		query,
		args)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, "Success")
}

// Create godoc
// @Summary      Update card by id
// @Description
// @Tags         cards
// @Accept       json
// @Produce      json
// @Param details body model.Card true "card Object to save"
// @Success      200  {object}  model.Card
// @Failure      400  {object}  error
// @Failure      500  {object}  error
// @Router       /admin/card/ [put]
func UpdateCard(c echo.Context) error {
	details := new(model.Card)
	if err := c.Bind(details); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	args := parseCard(*details)

	query := "UPDATE cards SET value=@value, suit=@suit, art=@art where id = @id"

	db := model.Pool()
	defer db.Close()

	_, err := db.Exec(
		context.Background(),
		query,
		args)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, "Card updated")
}

// Create godoc
// @Summary      Get card by barcode
// @Description
// @Tags         cards
// @Accept       json
// @Produce      json
// @Param        id    query     string  true  "search for card by id"'
// @Success      200  {object}  model.Card
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /player/card/ [get]
func GetCard(c echo.Context) error {
	id := c.Request().URL.Query().Get("id")

	db := model.Pool()
	defer db.Close()

	rows, err := db.Query(context.Background(), "SELECT * FROM cards WHERE id=$1", id)

	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}

	v := []model.Card{}

	for rows.Next() {
		var card model.Card

		err := rows.Scan(&card.Id, &card.Value, &card.Suit, &card.Art)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		v = append(v, card)
	}

	if len(v) > 1 {
		return c.JSON(http.StatusInternalServerError, "Too many cards with the same id found")
	}

	if len(v) == 0 {
		return c.JSON(http.StatusNotFound, "No card with ID found")
	}

	r, _ := json.Marshal(v[0])

	return c.JSON(http.StatusOK, string(r))
}

// Create godoc
// @Summary      Get all cards
// @Description
// @Tags         cards
// @Accept       json
// @Produce      json
// @Success      200  {object}  []model.GameplayCard{}
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /player/allcards/ [get]
func GetAllCards(c echo.Context) error {
	db := model.Pool()
	defer db.Close()

	rows, err := db.Query(context.Background(), "SELECT * FROM cards")

	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}

	v := []model.Card{}

	for rows.Next() {
		var card model.Card

		err := rows.Scan(&card.Id, &card.Value, &card.Suit, &card.Art)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		v = append(v, card)
	}

	if len(v) == 0 {
		return c.JSON(http.StatusNotFound, "No cards found. Something is very wrong")
	}

	r, _ := json.Marshal(v)

	return c.JSON(http.StatusOK, string(r))
}

// Create godoc
// @Summary      Get cards by ids
// @Description
// @Tags         cards
// @Accept       json
// @Produce      json
// @Param        ids    query     string  true  "csv of ids"'
// @Success      200  {object}  []model.GameplayCard{}
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /player/gameplaycards/ [get]
func GetGameplayCards(c echo.Context) error {
	ids := c.Request().URL.Query().Get("ids")
	v, err := QueryForCards(ids)
	if err != nil {
		return err
	}

	r, _ := json.Marshal(v)
	return c.JSON(http.StatusOK, string(r))
}

func QueryForCards(ids string) ([]model.GameplayCard, error) {
	db := model.Pool()
	defer db.Close()

	rows, err := db.Query(context.Background(), "SELECT * FROM gameplaycards NATURAL JOIN cards WHERE gameplaycards.id IN ("+ids+")")

	if err != nil {
		return []model.GameplayCard{}, err
	}

	v := []model.GameplayCard{}

	for rows.Next() {
		var card model.GameplayCard

		err := rows.Scan(&card.Id, &card.CardId, &card.OrigOwner, &card.CurrOwner, &card.State, &card.Value, &card.Suit, &card.Art)
		if err != nil {
			return []model.GameplayCard{}, err
		}

		v = append(v, card)
	}

	if len(v) == 0 {
		return []model.GameplayCard{}, errors.New("No cards found.")
	}

	return v, nil
}

// Create godoc
// @Summary      Get card by barcode
// @Description
// @Tags         cards
// @Accept       json
// @Produce      json
// @Param        barcode    query     string  true  "search for card by barcode"'
// @Param        storeId    query     string  true  "Store in which the barcode was found"'
// @Success      200  {object}  model.Card
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /card/ [delete]
func DeleteCard(c echo.Context) error {

	return c.JSON(http.StatusOK, nil)
}

func parseCard(details model.Card) pgx.NamedArgs {
	return pgx.NamedArgs{
		"id":    details.Id,
		"value": details.Value,
		"suit":  details.Suit,
		"art":   details.Art,
	}
}
