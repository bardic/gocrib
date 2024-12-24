package controller

import (
	"context"
	"slices"
	"time"

	"github.com/bardic/gocrib/queries/queries"

	conn "github.com/bardic/gocrib/server/db"
)

// UpdatePlayerById updates a player by id
func UpdatePlayerById(player *queries.Player) (queries.Player, error) {
	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	p, err := q.UpdatePlayer(ctx, queries.UpdatePlayerParams{
		Hand:    player.Hand,
		Play:    player.Play,
		Kitty:   player.Kitty,
		Score:   player.Score,
		Isready: player.Isready,
		Art:     player.Art,
		ID:      player.ID,
	})

	if err != nil {
		return queries.Player{}, err
	}

	return p, nil
}

// GetPlayerById gets a player by id
func GetPlayerById(id int32) (queries.Player, error) {
	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	p, err := q.GetPlayer(ctx, id)

	if err != nil {
		return queries.Player{}, err
	}

	return p, nil
}

// Equals compares two players
func Equals(p1, p2 queries.Player) bool {
	if p1.ID != p2.ID {
		return false
	}

	if slices.Equal(p1.Hand, p2.Hand) {
		return false
	}

	if slices.Equal(p1.Play, p2.Play) {
		return false
	}

	if slices.Equal(p1.Kitty, p2.Kitty) {
		return false
	}

	if p1.Score != p2.Score {
		return false
	}

	if p1.Isready != p2.Isready {
		return false
	}

	if p1.Art != p2.Art {
		return false
	}

	return true
}
