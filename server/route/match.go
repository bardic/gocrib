package route

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/bardic/cribbage/server/model"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

// Create godoc
// @Summary      Create new match
// @Description
// @Tags         match
// @Accept       json
// @Produce      json
// @Param details body model.Match true "match Object to save"
// @Success      200  {object}  model.Match
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /player/match/ [post]
func NewMatch(c echo.Context) error {
	details := new(model.Match)
	fmt.Print(time.Now())
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	args := parseMatch(*details)

	query := "INSERT INTO match(lobbyId, deckId, currentPlayerTurn, art) VALUES (@lobbyId, @deckId, @currentPlayerTurn, @art)"

	db := model.Pool()
	defer db.Close()

	_, err := db.Exec(
		context.Background(),
		query,
		args)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, "meow")
}

// Create godoc
// @Summary      Update match by barcode
// @Description
// @Tags         match
// @Accept       json
// @Produce      json
// @Param details body model.Match true "match Object to save"
// @Success      200  {object}  model.Match
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /player/match/ [put]
func UpdateMatch(c echo.Context) error {
	details := new(model.Match)
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	args := parseMatch(*details)

	if err := updateMatchQuery(args); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "meow")
}

func updateMatchQuery(args pgx.NamedArgs) error {
	query := "UPDATE match SET lobbyId = @lobbyId, deckId = @deckId, cardsInPlay = @cardsInPlay, cutGameCardId = @cutGameCardId,currentPlayerTurn = @currentPlayerTurn, turnPassTimestamps=@turnPassTimestamps, art=@art where id=@id"

	db := model.Pool()
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
// @Summary      Get match by barcode
// @Description
// @Tags         match
// @Accept       json
// @Produce      json
// @Param        id    query     string  true  "search for match by barcode"'
// @Success      200  {object}  model.Match
// @Failure      404  {object}  error
// @Failure      422  {object}  error
// @Router       /player/match/ [get]
func GetMatch(c echo.Context) error {
	p := c.Request().URL.Query().Get("id")
	id, err := strconv.Atoi(p)

	if err != nil {
		return err
	}

	v, err := GetMatchQuery(id)
	if err != nil {
		return err
	}

	r, _ := json.Marshal(v)

	return c.JSON(http.StatusOK, string(r))
}

func GetMatchQuery(id int) (model.Match, error) {
	db := model.Pool()
	defer db.Close()

	rows, err := db.Query(context.Background(), "SELECT json_agg( json_build_object('lobbyid', lobbyid, 'deckid', deckid, 'cardsinplay', cardsinplay, 'playerids', playerids, 'cutgamecardid', cutgamecardid, 'currentplayerturn', currentplayerturn, 'turnpasstimestamps', turnpasstimestamps, 'art', art,  'players', (SELECT json_agg(json_build_object( 'id', p.id,'hand', p.hand, 'kitty', p.kitty, 'score', p.score, 'art', p.art )) FROM player as p WHERE p.id = ANY(m.playerids)))) FROM match as m WHERE m.lobbyId = $1", id)
	if err != nil {
		return model.Match{}, err
	}

	v := []model.Match{}

	for rows.Next() {
		var match []model.Match

		err := rows.Scan(&match)
		if err != nil {
			return model.Match{}, err
		}

		v = append(v, match...)
	}

	if len(v) == 0 {
		return model.Match{}, errors.New("no match found")
	}

	return v[0], nil
}

// Create godoc
// @Summary      Get match by barcode
// @Description
// @Tags         match
// @Accept       json
// @Produce      json
// @Param        id    query     string  true  "search for match by barcode"'
// @Success      200  {object}  model.Match
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /admin/match/ [delete]
func DeleteMatch(c echo.Context) error {
	// b := c.Request().URL.Query().Get("barcode")
	// s := c.Request().URL.Query().Get("storeId")

	// match, _ := getmatch(b, s)

	return c.JSON(http.StatusOK, nil)
}

func UpdateCut(matchId int, cutCardId int) error {
	args := pgx.NamedArgs{"id": matchId, "cardId": cutCardId}

	query := "UPDATE match SET cutGameCardId = @cardId where id=@id"

	db := model.Pool()
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

func UpdateCardsInPlay(matchId int, gameplayCardId int) (model.Match, error) {
	args := pgx.NamedArgs{"id": matchId, "cardId": gameplayCardId}

	query := "UPDATE match SET cardsInPlay = array_append(cardsInPlay, @cardId)	where id=@id RETURNING *"

	db := model.Pool()
	defer db.Close()

	row := db.QueryRow(
		context.Background(),
		query,
		args)

	var match model.Match
	err := row.Scan(
		&match.Id,
		&match.LobbyId,
		&match.DeckId,
		&match.CardsInPlay,
		&match.PlayerIds,
		&match.CutGameCardId,
		&match.CurrentPlayerTurn,
		&match.TurnPassTimestamps,
		&match.Art,
	)
	if err != nil {
		fmt.Println(err)
		return model.Match{}, err
	}

	return match, nil
}

func parseMatch(details model.Match) pgx.NamedArgs {
	return pgx.NamedArgs{
		"id":                 details.Id,
		"lobbyId":            details.LobbyId,
		"deckId":             details.DeckId,
		"cardsInPlay":        details.CardsInPlay,
		"cutGameCardId":      details.CutGameCardId,
		"currentPlayerTurn":  details.CurrentPlayerTurn,
		"turnPassTimestamps": details.TurnPassTimestamps,
		"art":                details.Art,
	}
}
