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
