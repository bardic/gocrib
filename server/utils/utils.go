package utils

import (
	"context"

	conn "github.com/bardic/cribbage/server/db"

	"github.com/bardic/gocrib/model"
	"github.com/jackc/pgx/v5"
)

func QueryForCards(ids string) ([]model.GameplayCard, error) {
	db := conn.Pool()
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

	return v, nil
}

func UpdatePlay(details model.HandModifier) (*model.GameMatch, error) {
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
		return nil, err
	}

	return PlayCard(details)
}

func PlayCard(details model.HandModifier) (*model.GameMatch, error) {
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
		return nil, err
	}

	m, err := GetMatches(details.MatchId)

	if err != nil {
		return nil, err
	}

	return &m[0], nil
}

func PlayersReady(players []model.Player) bool {
	ready := true

	if len(players) < 2 {
		return false
	}

	for _, p := range players {
		if !p.IsReady {
			ready = false
		}
	}

	return ready
}

func GetMatchForPlayerId(playerId int) (int, error) {
	args := pgx.NamedArgs{
		"id": playerId,
	}

	query := `SELECT id from match WHERE @id = ANY(playerids)`

	var matchId int

	db := conn.Pool()
	defer db.Close()

	err := db.QueryRow(
		context.Background(),
		query,
		args).Scan(&matchId)

	if err != nil {
		return 0, err
	}

	return matchId, nil
}

func ParsePlayer(details model.Player) pgx.NamedArgs {
	return pgx.NamedArgs{
		"id":      details.Id,
		"play":    details.Play,
		"hand":    details.Hand,
		"kitty":   details.Kitty,
		"score":   details.Score,
		"isReady": details.IsReady,
		"art":     details.Art,
	}
}

func Deal(match *model.GameMatch) (*model.GameDeck, error) {
	deck, err := GetDeckById(match.DeckId)

	if err != nil {
		return nil, err
	}

	players := []model.Player{}

	for _, ids := range match.PlayerIds {
		player, err := GetPlayerById(ids)

		if err != nil {
			return nil, err
		}

		players = append(players, player)
	}

	cardsPerHand := 6
	if len(players) == 3 {
		cardsPerHand = 5
	}

	for i := 0; i < len(players)*cardsPerHand; i++ {
		var card model.GameplayCard
		card, deck.Cards = deck.Cards[0], deck.Cards[1:]
		idx := len(players) - 1 - i%len(players)

		if len(players[idx].Hand) < cardsPerHand {
			players[idx].Hand = append(players[idx].Hand, card.CardId)
		}
	}

	for _, p := range players {
		UpdatePlayerById(p)
	}

	return &deck, nil
}

func GetPlayerById(id int) (model.Player, error) {
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
		&player.IsReady,
		&player.Art,
	)

	if err != nil {
		return model.Player{}, err
	}

	return player, nil

}

func UpdatePlayerById(player model.Player) (model.Player, error) {
	args := ParsePlayer(player)

	query := `UPDATE player SET 
		hand = @hand, 
		play = @play, 
		kitty = @kitty, 
		score = @score, 
		isReady = @isReady,
		art = @art 
	where 
		id = @id`

	db := conn.Pool()
	defer db.Close()

	_, err := db.Exec(
		context.Background(),
		query,
		args)

	if err != nil {
		return model.Player{}, err
	}

	return player, nil
}
