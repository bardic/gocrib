package chat

import (
	"context"
	"fmt"
	"net/http"

	conn "github.com/bardic/cribbage/server/db"
	"github.com/bardic/cribbagev2/model"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

// Create godoc
// @Summary      Create new chat
// @Description
// @Tags         chats
// @Accept       json
// @Produce      json
// @Param details body model.Chat true "chat Object to save"
// @Success      200  {object}  model.Chat
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /chat/ [post]
func NewChat(c echo.Context) error {
	details := new(model.Chat)
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	fmt.Println(details)

	args := parseChat(*details)

	// "value":         details.Value,
	// 		"suit":          details.Suit,
	// 		"currentOwner":  details.Curr_owner,
	// 		"originalOwner": details.Orig_owner,
	// 		"state":         details.State,
	// 		"art":           details.Art,

	query := "INSERT INTO chats(value, suit, currentOwner, originalOwner, state, art) VALUES (@value, @suit, @currentOwner, @originalOwner, @state, @art)"

	fmt.Println("What the hell is happening")
	fmt.Println(query)

	db := conn.Pool()
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
// @Summary      Update chat by barcode
// @Description
// @Tags         chats
// @Accept       json
// @Produce      json
// @Param details body model.Chat true "chat Object to save"
// @Success      200  {object}  model.Chat
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /chat/ [put]
func UpdateChat(c echo.Context) error {
	details := new(model.Chat)
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	args := parseChat(*details)

	query := "UPDATE chats SET name = @name, cost = @cost, weight = @weight, unit=@unit, storename=@storeName, storeNeighborhood=@storeNeighborhood, tags=@tags where barcode = @barcode AND storeId = @storeId"

	db := conn.Pool()
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
// @Summary      Get chat by barcode
// @Description
// @Tags         chats
// @Accept       json
// @Produce      json
// @Param        barcode    query     string  true  "search for chat by barcode"'
// @Param        storeId    query     string  true  "Store in which the barcode was found"'
// @Success      200  {object}  model.Chat
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /chat/ [get]
func GetChat(c echo.Context) error {
	id := c.Request().URL.Query().Get("id")

	db := conn.Pool()
	defer db.Close()

	rows, err := db.Query(context.Background(), "SELECT * FROM chat WHERE id=$1", id)

	v := []model.Chat{}

	for rows.Next() {
		var card model.Chat

		err := rows.Scan(&card.Id, &card.Members, &card.Messages)
		if err != nil {
			return err
		}

		v = append(v, card)
	}

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, v)
}

// Create godoc
// @Summary      Get chat by barcode
// @Description
// @Tags         chats
// @Accept       json
// @Produce      json
// @Param        barcode    query     string  true  "search for chat by barcode"'
// @Param        storeId    query     string  true  "Store in which the barcode was found"'
// @Success      200  {object}  model.Chat
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /chat/ [delete]
func DeleteChat(c echo.Context) error {
	// b := c.Request().URL.Query().Get("barcode")
	// s := c.Request().URL.Query().Get("storeId")

	// chat, _ := getchat(b, s)

	return c.JSON(http.StatusOK, nil)
}

func parseChat(details model.Chat) pgx.NamedArgs {
	return pgx.NamedArgs{
		"Members":  details.Members,
		"Messages": details.Messages,
	}
}
