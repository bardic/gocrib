package route

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	conn "github.com/bardic/cribbage/server/db"
	"github.com/bardic/cribbagev2/model"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

// Create godoc
// @Summary      Create new match
// @Description
// @Tags         match
// @Accept       json
// @Produce      json
// @Param details body model.MatchRequirements true "MatchRequirements"
// @Success      200  {object}  model.Match
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /player/match/ [post]
func NewMatch(c echo.Context) error {
	details := new(model.MatchRequirements)
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	db := conn.Pool()
	defer db.Close()

	p1 := "1"

	d, err := newDeck()

	if err != nil {
		return err
	}

	match := model.GameMatch{}
	match.DeckId = d.Id

	p1Id, err := strconv.Atoi(p1)

	if err != nil {
		return err
	}

	match.PlayerIds = []int{p1Id}
	match.EloRangeMin = details.EloRangeMin
	match.EloRangeMax = details.EloRangeMax
	match.PrivateMatch = details.IsPrivate

	args := parseMatch(match)

	query := `INSERT INTO match(
				playerIds,
				privateMatch,
				eloRangeMin,
				eloRangeMax,
				deckId,
				cutGameCardId,
				currentPlayerTurn,
				turnPassTimestamps,
				gameState,
				art)
			VALUES (
				@playerIds,
				@privateMatch,
				@eloRangeMin,
				@eloRangeMax,
				@deckId,
				@cutGameCardId,
				@currentPlayerTurn,
				@turnPassTimestamps,
				@gameState,
				@art)
			RETURNING id`

	var matchId int
	err = db.QueryRow(
		context.Background(),
		query,
		args).Scan(&matchId)

	if err != nil {
		return err
	}

	// d = *d.Shuffle()
	// for i := 0; i < 12; i++ {
	// 	if i%2 == 0 {
	// 		match.Players[0].Hand = append(match.Players[0].Hand, d.Cards[i].CardId)
	// 	} else {
	// 		match.Players[1].Hand = append(match.Players[1].Hand, d.Cards[i].CardId)
	// 	}
	// }

	// match.GameState = model.DiscardState

	// // err = updateMatch(match)

	// if err != nil {
	// 	return err
	// }

	// if match.Players[0].Id == details.RequesterId {
	// 	match.Players[1].Hand = []int{}
	// } else {
	// 	match.Players[0].Hand = []int{}
	// }

	match.Id = matchId
	match.GameState = model.WaitingState
	return c.JSON(http.StatusOK, match)
}

func findMatchInEloRange(req model.MatchRequirements) (model.GameMatch, error) {
	db := conn.Pool()
	defer db.Close()

	args := pgx.NamedArgs{"eloMin": req.EloRangeMin, "eloMax": req.EloRangeMax}
	q := `SELECT
			json_build_object(
				'id', id,
				'playerIds', playerIds,
				'creationDate', creationDate,
				'privateMatch', privateMatch,
				'eloRangeMin', eloRangeMin,
				'eloRangeMax', eloRangeMax,
				'deckid', deckid,
				'cutgamecardid', cutgamecardid,
				'currentplayerturn', currentplayerturn,
				'turnpasstimestamps', turnpasstimestamps,
				'art', art,
				'players', (SELECT json_agg(
					json_build_object(
						'id', p.id,
						'play', p.play,
						'hand', p.hand,
						'kitty', p.kitty,
						'score', p.score,
						'art', p.art ))
				FROM player as p WHERE p.id = ANY(m.playerIds)))
			FROM match as m WHERE m.eloRangeMin BETWEEN @eloMin AND @eloMax OR m.eloRangeMax BETWEEN @eloMin AND @eloMax`

	var match model.GameMatch
	err := db.QueryRow(
		context.Background(),
		q,
		args).Scan(&match)

	if err != nil && err != pgx.ErrNoRows {
		return model.GameMatch{}, err
	}

	return match, nil
}

// Create godoc
// @Summary      Join match by id
// @Description
// @Tags         match
// @Accept       json
// @Produce      json
// @Param details body model.JoinMatchReq true "match Object to update"
// @Success      200  {object}  model.Match
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /player/match/join [put]
func JoinMatch(c echo.Context) error {
	details := new(model.JoinMatchReq)
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := updatePlayersInMatch(*details); err != nil {
		return err
	}

	empty, err := json.Marshal(model.Match{})

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, string(empty))
}

func updatePlayersInMatch(req model.JoinMatchReq) error {
	args := pgx.NamedArgs{
		"matchId":     req.MatchId,
		"requesterId": req.RequesterId,
		"gameState":   model.WaitingState,
	}
	query := `UPDATE match SET
				playerIds=ARRAY_APPEND(playerIds, @requesterId),
				gameState=@gameState
			where id=@matchId`

	db := conn.Pool()
	defer db.Close()

	_, err := db.Exec(
		context.Background(),
		query,
		args)

	if err != nil {
		return err
	}

	m, err := getMatch(req.MatchId)

	if err != nil {
		return err
	}

	deal(m)

	return nil
}

// Create godoc
// @Summary      Update match by id
// @Description
// @Tags         match
// @Accept       json
// @Produce      json
// @Param details body model.Match true "match Object to update"
// @Success      200  {object}  model.Match
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /player/match/ [put]
func UpdateMatch(c echo.Context) error {
	details := new(model.GameMatch)
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := updateMatch(*details); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "meow")
}

func updateMatch(match model.GameMatch) error {
	args := parseMatch(match)
	query := `UPDATE match SET
				playerIds = @playerIds,
				creationDate = @creationDate,
				privateMatch = @privateMatch,
				eloRangeMin = @eloRangeMin,
				eloRangeMax = @eloRangeMax,
				deckId = @deckId,
				cutGameCardId = @cutGameCardId,
				currentPlayerTurn = @currentPlayerTurn,
				turnPassTimestamps = @turnPassTimestamps,
				gameState= @gameState,
				art = @art
			where id=@id`

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

func updateGameState(matchId int, state model.GameState) error {
	args := pgx.NamedArgs{
		"id":        matchId,
		"gameState": state,
	}
	query := `UPDATE match SET
				gameState=@gameState
			where id=@id`

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
// @Summary      Get match by id
// @Description
// @Tags         match
// @Accept       json
// @Produce      json
// @Param        id    query     string  true  "search for match by id"'
// @Success      200  {object}  []model.Match
// @Failure      404  {object}  error
// @Failure      422  {object}  error
// @Router       /player/matches/ [get]
func GetMatches(c echo.Context) error {
	p := c.Request().URL.Query().Get("id")
	id, err := strconv.Atoi(p)

	if err != nil {
		return err
	}

	v, err := getMatches(id)
	if err != nil {
		return err
	}

	r, _ := json.Marshal(v)

	return c.JSON(http.StatusOK, string(r))
}

// Create godoc
// @Summary      Get match by id
// @Description
// @Tags         match
// @Accept       json
// @Produce      json
// @Success      200  {object}  []model.Match
// @Failure      404  {object}  error
// @Failure      422  {object}  error
// @Router       /player/matches/open [get]
func GetOpenMatches(c echo.Context) error {
	v, err := getOpenMatches()
	if err != nil {
		return err
	}

	r, _ := json.Marshal(v)

	return c.JSON(http.StatusOK, string(r))
}

// Create godoc
// @Summary      Get match by id
// @Description
// @Tags         match
// @Accept       json
// @Produce      json
// @Param        id    query     string  true  "search for match by id"'
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

	v, err := getMatch(id)
	if err != nil {
		return err
	}

	r, _ := json.Marshal(v)

	return c.JSON(http.StatusOK, string(r))
}

// Create godoc
// @Summary      Get deck by id
// @Description
// @Tags         deck
// @Accept       json
// @Produce      json
// @Param        id    query     string  true  "search for deck by id"'
// @Success      200  {object}  model.GameDeck
// @Failure      404  {object}  error
// @Failure      422  {object}  error
// @Router       /player/match/deck/ [get]
func GetDeck(c echo.Context) error {
	p := c.Request().URL.Query().Get("id")
	id, err := strconv.Atoi(p)

	if err != nil {
		return err
	}

	deck, err := getDeck(id)

	if err != nil {
		return err
	}

	r, err := json.Marshal(deck)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, string(r))
}

func getDeck(id int) (model.GameDeck, error) {
	db := conn.Pool()
	defer db.Close()
	var deckId int
	var cards []model.GameplayCard
	err := db.QueryRow(context.Background(), "SELECT * FROM deck WHERE id=$1", id).Scan(&deckId, &cards)

	if err != nil {
		return model.GameDeck{}, err
	}

	deck := model.GameDeck{
		Id:    deckId,
		Cards: cards,
	}

	return deck, nil
}

func getMatches(id int) ([]model.Match, error) {
	db := conn.Pool()
	defer db.Close()

	var matches []model.Match
	rows, err := db.Query(context.Background(), `SELECT
			json_build_object(
				'id', id,
				'playerIds', playerIds,
				'creationDate', creationDate,
				'privateMatch', privateMatch,
				'eloRangeMin', eloRangeMin,
				'eloRangeMax', eloRangeMax,
				'deckid', deckid,
				'cutgamecardid', cutgamecardid,
				'currentplayerturn', currentplayerturn,
				'turnpasstimestamps', turnpasstimestamps,
				'gamestate', gamestate,
				'art', art,
				'players', (SELECT json_agg(
					json_build_object(
						'id', p.id,
						'play', p.play,
						'hand', p.hand,
						'kitty', p.kitty,
						'score', p.score,
						'art', p.art ))
				FROM player as p WHERE p.id = ANY(m.playerIds)))
			FROM match as m WHERE $1=ANY(m.playerIds)`, id)

	if err != nil {
		return []model.Match{}, err
	}

	for rows.Next() {
		var match model.Match

		err := rows.Scan(&match)
		if err != nil {
			fmt.Println(err)
			return []model.Match{}, &echo.BindingError{}
		}

		matches = append(matches, match)
	}

	return matches, nil
}

func getMatch(id int) (model.Match, error) {
	db := conn.Pool()
	defer db.Close()

	var match model.Match
	err := db.QueryRow(context.Background(), `SELECT
			json_build_object(
				'id', id,
				'playerIds', playerIds,
				'creationDate', creationDate,
				'privateMatch', privateMatch,
				'eloRangeMin', eloRangeMin,
				'eloRangeMax', eloRangeMax,
				'deckid', deckid,
				'cutgamecardid', cutgamecardid,
				'currentplayerturn', currentplayerturn,
				'turnpasstimestamps', turnpasstimestamps,
				'gamestate', gamestate,
				'art', art,
				'players', (SELECT json_agg(
					json_build_object(
						'id', p.id,
						'play', p.play,
						'hand', p.hand,
						'kitty', p.kitty,
						'score', p.score,
						'art', p.art ))
				FROM player as p WHERE p.id = ANY(m.playerIds)))
			FROM match as m WHERE m.id = $1`, id).Scan(
		&match,
	)

	if err != nil {
		return model.Match{}, err
	}

	return match, nil
}

func getOpenMatches() ([]model.Match, error) {
	db := conn.Pool()
	defer db.Close()

	rows, err := db.Query(context.Background(), `SELECT
			json_build_object(
				'id', id,
				'playerIds', playerIds,
				'creationDate', creationDate,
				'privateMatch', privateMatch,
				'eloRangeMin', eloRangeMin,
				'eloRangeMax', eloRangeMax,
				'deckid', deckid,
				'cutgamecardid', cutgamecardid,
				'currentplayerturn', currentplayerturn,
				'turnpasstimestamps', turnpasstimestamps,
				'gamestate', gamestate,
				'art', art,
				'players', (SELECT json_agg(
					json_build_object(
						'id', p.id,
						'play', p.play,
						'hand', p.hand,
						'kitty', p.kitty,
						'score', p.score,
						'art', p.art ))
				FROM player as p WHERE p.id = ANY(m.playerIds)))
			FROM match as m`)

	var matches []model.Match

	for rows.Next() {
		var match model.Match

		err := rows.Scan(&match)
		if err != nil {
			return []model.Match{}, err
		}

		matches = append(matches, match)
	}

	if err != nil && err != pgx.ErrNoRows {
		return []model.Match{}, err
	}

	return matches, nil
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

func UpdateCardsInPlay(a *model.GameAction) (model.Match, error) {
	args := pgx.NamedArgs{"id": a.MatchId, "cardsIds": a.CardsIds}

	query := "UPDATE match SET cardsInPlay = array_append(cardsInPlay, @cardIds) where id=@id RETURNING *"

	db := conn.Pool()
	defer db.Close()

	var match model.Match
	err := db.QueryRow(
		context.Background(),
		query,
		args).Scan(
		&match.Id,
		&match.PlayerIds,
		&match.CreationDate,
		&match.PrivateMatch,
		&match.EloRangeMin,
		&match.EloRangeMax,
		&match.DeckId,
		&match.CutGameCardId,
		&match.CurrentPlayerTurn,
		&match.TurnPassTimestamps,
		&match.GameState,
		&match.Art,
	)

	if err != nil {
		return model.Match{}, err
	}

	return match, nil
}

func parseMatch(details model.GameMatch) pgx.NamedArgs {
	return pgx.NamedArgs{
		"id":                 details.Id,
		"gameState":          details.GameState,
		"playerIds":          details.PlayerIds,
		"privateMatch":       details.PrivateMatch,
		"eloRangeMin":        details.EloRangeMin,
		"eloRangeMax":        details.EloRangeMax,
		"deckId":             details.DeckId,
		"creationDate":       details.CreationDate,
		"cutGameCardId":      details.CutGameCardId,
		"currentPlayerTurn":  details.CurrentPlayerTurn,
		"turnPassTimestamps": []string{},
		"art":                details.Art,
	}
}

func newDeck() (model.GameDeck, error) {
	db := conn.Pool()
	defer db.Close()

	rows, err := db.Query(context.Background(), "SELECT * FROM cards")

	v := []model.Card{}

	for rows.Next() {
		var card model.Card

		err := rows.Scan(&card.Id, &card.Value, &card.Suit, &card.Art)
		if err != nil {
			return model.GameDeck{}, err
		}

		v = append(v, card)
	}

	if err != nil {
		return model.GameDeck{}, err
	}

	deck := model.GameDeck{
		Cards: []model.GameplayCard{},
	}

	for _, c := range v {
		deck.Cards = append(deck.Cards, model.GameplayCard{
			CardId: c.Id,
			State:  0,
		})
	}

	b, err := json.Marshal(deck.Cards)

	if err != nil {
		return model.GameDeck{}, err
	}

	args := pgx.NamedArgs{"cards": string(b)}

	query := "INSERT INTO deck(cards) VALUES (@cards) RETURNING id"

	var deckId int
	err = db.QueryRow(
		context.Background(),
		query,
		args).Scan(&deckId)

	if err != nil {
		return model.GameDeck{}, err
	}

	deck.Id = deckId

	return deck, nil
}
