package route

import (
	"context"
	"encoding/json"
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

	args := pgx.NamedArgs{"eloMin": details.EloRangeMin, "eloMax": details.EloRangeMax}
	q := `SELECT 
			json_build_object(
				'playerIds', playerIds, 
				'creationDate', creationDate,
				'privateMatch', privateMatch,
				'eloRangeMin', eloRangeMin,
				'eloRangeMax', eloRangeMax,
				'deckid', deckid, 
				'cardsinplay', cardsinplay, 
				'cutgamecardid', cutgamecardid, 
				'currentplayerturn', currentplayerturn, 
				'turnpasstimestamps', turnpasstimestamps, 
				'art', art,  
				'players', (SELECT json_agg(
					json_build_object( 
						'id', p.id,
						'hand', p.hand, 
						'kitty', p.kitty, 
						'score', p.score, 
						'art', p.art )) 
				FROM player as p WHERE p.id = ANY(m.playerIds))) 
			FROM match as m WHERE m.eloRangeMin BETWEEN @eloMin AND @eloMax OR m.eloRangeMax BETWEEN @eloMin AND @eloMax`

	var match model.Match
	err := db.QueryRow(
		context.Background(),
		q,
		args).Scan(&match)

	if err != nil && err != pgx.ErrNoRows {
		return err
	}

	match.EloRangeMin = details.EloRangeMin
	match.EloRangeMax = details.EloRangeMax
	match.PrivateMatch = details.IsPrivate

	p1, err := newPlayer()
	if err != nil {
		return err
	}

	p2, err := newPlayer()

	if err != nil {
		return err
	}

	d, err := newDeck()

	if err != nil {
		return err
	}

	match.DeckId = d.Id
	match.PlayerIds = []int{p1, p2}

	args = parseMatch(match)

	query := `INSERT INTO match( 	
				playerIds,
				privateMatch,
				eloRangeMin,
				eloRangeMax,
				deckId,
				currentPlayerTurn,
				gameState, 
				art) 
			VALUES ( 
				@playerIds,
				@privateMatch,
				@eloRangeMin,
				@eloRangeMax,
				@deckId,
				@currentPlayerTurn, 
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

	match.Id = matchId

	return c.JSON(http.StatusOK, match)
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
	details := new(model.Match)
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	args := parseMatch(*details)

	if err := updateMatch(args); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "meow")
}

func updateMatch(args pgx.NamedArgs) error {
	query := `UPDATE match SET  
				playerIds = @playerIds
				creationDate = @creationDate,
				privateMatch = @privateMatch,
				eloRangeMin = @eloRangeMin,
				eloRangeMax = @eloRangeMax,
				deckId = @deckId,
				cardsInPlay = @cardsInPlay,
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
// @Router       /match/deck/ [get]
func GetDeck(c echo.Context) error {
	p := c.Request().URL.Query().Get("id")
	id, err := strconv.Atoi(p)

	if err != nil {
		return err
	}

	db := conn.Pool()
	defer db.Close()
	var deckId int
	var cards []model.GameplayCard
	err = db.QueryRow(context.Background(), "SELECT * FROM deck WHERE id=$1", id).Scan(&deckId, &cards)

	if err != nil {
		return err
	}

	r, err := json.Marshal(cards)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, string(r))
}

func getMatch(id int) (model.Match, error) {
	db := conn.Pool()
	defer db.Close()

	var match model.Match
	err := db.QueryRow(context.Background(), `SELECT 
			json_build_object(
				'playerIds', playerIds, 
				'creationDate', creationDate,
				'privateMatch', privateMatch,
				'eloRangeMin', eloRangeMin,
				'eloRangeMax', eloRangeMax,
				'deckid', deckid, 
				'cardsinplay', cardsinplay, 
				'cutgamecardid', cutgamecardid, 
				'currentplayerturn', currentplayerturn, 
				'turnpasstimestamps', turnpasstimestamps, 
				'art', art,  
				'players', (SELECT json_agg(
					json_build_object( 
						'id', p.id,
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

func UpdateCardsInPlay(matchId int, gameplayCardId int) (model.Match, error) {
	args := pgx.NamedArgs{"id": matchId, "cardId": gameplayCardId}

	query := "UPDATE match SET cardsInPlay = array_append(cardsInPlay, @cardId)	where id=@id RETURNING *"

	db := conn.Pool()
	defer db.Close()

	var match model.Match
	err := db.QueryRow(
		context.Background(),
		query,
		args).Scan(
		&match.Id,
		&match.GameState,
		&match.PlayerIds,
		&match.PrivateMatch,
		&match.EloRangeMin,
		&match.EloRangeMax,
		&match.DeckId,
		&match.CardsInPlay,
		&match.CutGameCardId,
		&match.CurrentPlayerTurn,
		&match.TurnPassTimestamps,
		&match.Art,
	)

	if err != nil {
		return model.Match{}, err
	}

	return match, nil
}

func parseMatch(details model.Match) pgx.NamedArgs {
	return pgx.NamedArgs{
		"id":                 details.Id,
		"gameState":          details.GameState,
		"playerIds":          details.PlayerIds,
		"privateMatch":       details.PrivateMatch,
		"eloRangeMin":        details.EloRangeMin,
		"eloRangeMax":        details.EloRangeMax,
		"deckId":             details.DeckId,
		"cardsInPlay":        details.CardsInPlay,
		"creationDate":       details.CreationDate,
		"cutGameCardId":      details.CutGameCardId,
		"currentPlayerTurn":  details.CurrentPlayerTurn,
		"turnPassTimestamps": details.TurnPassTimestamps,
		"art":                details.Art,
	}
}

func newPlayer() (int, error) {
	args := pgx.NamedArgs{"art": "default.png"}
	query := "INSERT INTO player (art) VALUES (@art) RETURNING id"

	db := conn.Pool()
	defer db.Close()

	var playerId int
	err := db.QueryRow(
		context.Background(),
		query,
		args).Scan(&playerId)

	if err != nil {
		return 0, err
	}

	return playerId, nil
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
