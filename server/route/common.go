package route

import (
	"context"
	"encoding/json"
	"time"

	"github.com/bardic/gocrib/queries/queries"
	conn "github.com/bardic/gocrib/server/db"
	"github.com/bardic/gocrib/vo"
)

const (
	CtxTimeout = 5 * time.Second
)

func UpdateGameState(matchID int, state queries.Gamestate) error {
	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		CtxTimeout,
	)
	defer cancel()

	_, err := q.UpdateMatchState(ctx, queries.UpdateMatchStateParams{Gamestate: state, ID: matchID})
	if err != nil {
		return err
	}

	return nil
}

func GetMatch(id int) (*vo.Match, error) {
	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		CtxTimeout,
	)
	defer cancel()

	m, err := q.GetMatchById(ctx, id)
	if err != nil {
		return nil, err
	}

	var match *vo.Match
	err = json.Unmarshal(m, &match)
	if err != nil {
		return nil, err
	}

	// TODO Add player info here to the gamematch object

	// players, err := q.GetPlayersByMatchId(ctx, id)

	// if err != nil {
	// 	return nil, err
	// }

	// playersForMatch, err := q.GetPlayersForMatchId(ctx, id)
	// // var player vo.GamePlayer
	// // err = json.Unmarshal(playersForMatch, &player)

	// if err != nil {
	// 	return nil, err

	// }

	// var gamePlayers []vo.GamePlayer
	// for _, player := range playersForMatch {
	// 	var gameplayer vo.GamePlayer
	// 	err = json.Unmarshal(player, &gameplayer)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	gamePlayers = append(gamePlayers, gameplayer)
	// }

	// updatedGameplayers := []*vo.GamePlayer{}
	// for _, p := range gamePlayers {

	// 	player := &vo.GamePlayer{
	// 		Player: queries.Player{
	// 			ID:        p.ID,
	// 			Accountid: p.Accountid,
	// 			Score:     p.Score,
	// 			Isready:   p.Isready,
	// 			Art:       p.Art,
	// 		},
	// 		Hand:  p.Hand,
	// 		Play:  p.Play,
	// 		Kitty: p.Kitty,
	// 	}

	// 	updatedGameplayers = append(updatedGameplayers, player)
	// }

	// match.Players = updatedGameplayers

	return match, nil
}
