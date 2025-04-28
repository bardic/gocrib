package helpers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/bardic/gocrib/queries/queries"

	conn "github.com/bardic/gocrib/server/db"
	"github.com/bardic/gocrib/vo"
)

func UpdateGameState(matchId *int, state queries.Gamestate) error {
	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := q.UpdateMatchState(ctx, queries.UpdateMatchStateParams{Gamestate: state, ID: matchId})

	if err != nil {
		return err
	}

	return nil
}

func GetMatch(id *int) (*vo.GameMatch, error) {
	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	m, err := q.GetMatchById(ctx, id)

	if err != nil {
		return nil, err
	}

	var match *vo.GameMatch
	err = json.Unmarshal(m, &match)
	if err != nil {
		return nil, err
	}

	//TODO Add player info here to the gamematch object

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

func GetOpenMatches() ([]queries.Match, error) {
	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	matchesData, err := q.GetOpenMatches(ctx, queries.GamestateNew)

	if err != nil {
		return nil, err
	}

	// var matches []queries.Match

	// for _, matchData := range matchesData {
	// 	var match queries.Match
	// 	// err = json.Unmarshal(matchData, &match)
	// 	// if err != nil {
	// 	// 	return nil, err
	// 	// }

	// 	match = queries.Match{}

	// 	matches = append(matches, match)
	// }

	return matchesData, nil
}
