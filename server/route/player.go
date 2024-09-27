package route

import (
	"context"
	"encoding/json"
	"net/http"

	conn "github.com/bardic/cribbage/server/db"
	"github.com/bardic/cribbagev2/model"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

// Create godoc
// @Summary      Create new player
// @Description
// @Tags         players
// @Accept       json
// @Produce      json
// @Param details body model.Player true "player Object to save"
// @Success      200  {object}  model.Player
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /player/player/ [post]
func NewPlayer(c echo.Context) error {
	details := new(model.Player)
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	args := parsePlayer(*details)

	query := "INSERT INTO player(hand, kitty, score, art) VALUES (@hand, @kitty, @score, @art)"

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
// @Summary      Update player by barcode
// @Description
// @Tags         players
// @Accept       json
// @Produce      json
// @Param details body model.Player true "player Object to save"
// @Success      200  {object}  model.Player
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /player/player/ [put]
func UpdatePlayer(c echo.Context) error {
	details := new(model.Player)
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err := UpdatePlayerQuery(*details)

	if err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return c.JSON(http.StatusOK, nil)
}

func UpdatePlayerQuery(player model.Player) error {
	args := parsePlayer(player)

	query := "UPDATE player SET hand = @hand, play = @play, kitty = @kitty, score = @score, art = @art where id = @id"

	db := conn.Pool()
	defer db.Close()

	_, err := db.Exec(
		context.Background(),
		query,
		args)

	if err != nil {
		return err
	}

	return nil
}

// Create godoc
// @Summary      Get player by barcode
// @Description
// @Tags         players
// @Accept       json
// @Produce      json
// @Param        id    query     string  true  "search for match by barcode"'
// @Success      200  {object}  model.Player
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /player/player/ [get]
func GetPlayer(c echo.Context) error {
	id := c.Request().URL.Query().Get("id")

	p, err := getPlayer(id)

	if err != nil {
		return err
	}

	r, _ := json.Marshal(p)
	return c.JSON(http.StatusOK, string(r))
}

func getPlayer(id string) (model.Player, error) {
	db := conn.Pool()
	defer db.Close()

	player := model.Player{}
	err := db.QueryRow(context.Background(), "SELECT * FROM player WHERE id=$1", id).Scan(
		&player.Id,
		&player.AccountId,
		&player.Play,
		&player.Hand,
		&player.Kitty,
		&player.Score,
		&player.Art,
	)

	if err != nil {
		return model.Player{}, err
	}

	return player, nil

}

// Create godoc
// @Summary      Get player by barcode
// @Description
// @Tags         players
// @Accept       json
// @Produce      json
// @Param        id    query     string  true  "search for match by barcode"'
// @Success      200  {object}  model.Player
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /player/player/ [delete]
func DeletePlayer(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}

func parsePlayer(details model.Player) pgx.NamedArgs {
	return pgx.NamedArgs{
		"id":    details.Id,
		"play":  details.Play,
		"hand":  details.Hand,
		"kitty": details.Kitty,
		"score": details.Score,
		"art":   details.Art,
	}
}

// Create godoc
// @Summary      Update kitty with ids
// @Description
// @Tags         match
// @Accept       json
// @Produce      json
// @Param details body model.HandModifier true "array of ids to add to kitty"
// @Success      200  {object}  model.Match
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /player/kitty [put]
func UpdateKitty(c echo.Context) error {
	details := &model.HandModifier{}
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	m, err := updateKitty(*details)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	m.GameState = model.PlayState

	r, _ := json.Marshal(m)

	return c.JSON(http.StatusOK, string(r))
}

func updateKitty(details model.HandModifier) (model.Match, error) {
	db := conn.Pool()
	defer db.Close()

	args := pgx.NamedArgs{
		"matchId": details.MatchId,
	}

	q := "SELECT currentplayerturn FROM match WHERE id = @matchId"

	var dealerPlayerId int
	err := db.QueryRow(
		context.Background(),
		q,
		args).Scan(&dealerPlayerId)

	if err != nil {
		return model.Match{}, err
	}

	args = pgx.NamedArgs{
		"dealerId": 1,
		"kitty":    details.CardIds,
	}

	q = "UPDATE player SET kitty = kitty + @kitty where id = @dealerId"

	_, err = db.Exec(
		context.Background(),
		q,
		args)

	if err != nil {
		return model.Match{}, err
	}

	return playCard(details)
}

// Create godoc
// @Summary      Update play with ids
// @Description
// @Tags         match
// @Accept       json
// @Produce      json
// @Param details body model.HandModifier true "array of ids to add to play"
// @Success      200  {object}  model.Match
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /player/play [put]
func UpdatePlay(c echo.Context) error {
	details := &model.HandModifier{}
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	m, err := updatePlay(*details)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	m.GameState = model.PlayState

	r, _ := json.Marshal(m)

	return c.JSON(http.StatusOK, string(r))
}

func updatePlay(details model.HandModifier) (model.Match, error) {
	db := conn.Pool()
	defer db.Close()

	args := pgx.NamedArgs{
		"playerId": details.PlayerId,
		"play":     details.CardIds,
	}

	q := "UPDATE player SET play = play + @play where id = @playerId"

	_, err := db.Exec(
		context.Background(),
		q,
		args)

	if err != nil {
		return model.Match{}, err
	}

	return playCard(details)
}

func playCard(details model.HandModifier) (model.Match, error) {
	db := conn.Pool()
	defer db.Close()

	args := pgx.NamedArgs{
		"playerId": details.PlayerId,
		"cards":    details.CardIds,
	}

	q := "UPDATE player SET hand = hand - @cards where id = @playerId"

	_, err := db.Exec(
		context.Background(),
		q,
		args)

	if err != nil {
		return model.Match{}, err
	}

	m, err := getMatch(details.MatchId)

	if err != nil {
		return model.Match{}, err
	}

	return m, nil
}
