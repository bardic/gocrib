package utils

import (
	"context"
	"encoding/json"

	conn "github.com/bardic/cribbage/server/db"
	"github.com/bardic/gocrib/model"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

var matchQuery = `SELECT
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
				'gameState', gameState,
				'players', (SELECT json_agg(
					json_build_object(
						'id', p.id,
						'accountid', p.accountid,
						'play', p.play,
						'hand', p.hand,
						'kitty', p.kitty,
						'score', p.score,
						'art', p.art ))
				FROM player as p WHERE p.Id = ANY(m.playerIds)))`

func CardsInPlay(players []model.Player) []int {
	play := []int{}
	for _, player := range players {
		play = append(play, player.Play...)
	}

	return play
}

func findMatchInEloRange(req model.MatchRequirements) (model.GameMatch, error) {
	db := conn.Pool()
	defer db.Close()

	args := pgx.NamedArgs{"eloMin": req.EloRangeMin, "eloMax": req.EloRangeMax}
	q := matchQuery + `FROM match 
	AS m WHERE m.eloRangeMin BETWEEN @eloMin AND @eloMax 
	OR m.eloRangeMax BETWEEN @eloMin AND @eloMax`

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

func UpdateGameState(matchId int, state model.GameState) error {
	args := pgx.NamedArgs{
		"id":        matchId,
		"gameState": state,
	}
	query := `UPDATE match SET
				gameState=@gameState
			WHERE id=@id`

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

func GetMatches(id int) ([]model.GameMatch, error) {
	db := conn.Pool()
	defer db.Close()

	var matches []model.GameMatch
	rows, err := db.Query(
		context.Background(),
		matchQuery+`FROM match as m WHERE $1=ANY(m.playerIds)`,
		id,
	)

	if err != nil {
		return []model.GameMatch{}, err
	}

	for rows.Next() {
		var match model.GameMatch

		err := rows.Scan(&match)
		if err != nil {

			return []model.GameMatch{}, &echo.BindingError{}
		}

		matches = append(matches, match)
	}

	return matches, nil
}

func GetMatch(id int) (*model.GameMatch, error) {
	db := conn.Pool()
	defer db.Close()

	var match *model.GameMatch
	err := db.QueryRow(
		context.Background(),
		matchQuery+" FROM match as m WHERE m.id = $1",
		id,
	).Scan(
		&match,
	)

	if err != nil {
		return nil, err
	}

	return match, nil
}

func GetOpenMatches() ([]model.GameMatch, error) {
	db := conn.Pool()
	defer db.Close()

	rows, err := db.Query(
		context.Background(),
		matchQuery+`FROM match as m`,
	)

	if err != nil && err != pgx.ErrNoRows {
		return []model.GameMatch{}, err
	}

	matches := []model.GameMatch{}

	for rows.Next() {
		var match model.GameMatch

		err := rows.Scan(&match)
		if err != nil {
			return []model.GameMatch{}, err
		}

		matches = append(matches, match)
	}

	return matches, nil
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

func ParseMatch(details model.GameMatch) pgx.NamedArgs {
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
		"players":            details.Players,
	}
}

func NewDeck() (model.GameDeck, error) {
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
			Card:   c,
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

func IsMatchReadyToStart(m *model.GameMatch) (bool, error) {
	if len(m.Players) == 2 {
		return true, nil
	}

	return false, nil
}

func UpdateMatchState(matchId int, state model.GameState) error {
	args := pgx.NamedArgs{
		"id":        matchId,
		"gameState": state,
	}

	query := `UPDATE match SET
					gameState= @gameState
				WHERE id=@id`

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
