package player

import (
	"context"
	"time"

	"github.com/bardic/gocrib/queries/queries"

	conn "github.com/bardic/gocrib/server/db"
)

func NewPlayerQuery(matchId, accountId int32) (*queries.Player, error) {

	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	p, err := q.CreatePlayer(ctx, queries.CreatePlayerParams{
		Accountid: accountId,
		Hand:      []int32{},
		Kitty:     []int32{},
		Play:      []int32{},
		Score:     0,
		Isready:   false,
		Art:       "default.png",
	})

	if err != nil {
		return nil, err
	}

	return &p, nil
}
