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

	err := q.UpdateMatchState(ctx, queries.UpdateMatchStateParams{Gamestate: state, ID: matchId})

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

	b, err := q.GetPlayerJSON(ctx, id)

	if err != nil {
		return nil, err

	}

	var gamePlayers []vo.GamePlayer
	for _, gp := range b {
		var gameplayer vo.GamePlayer
		err = json.Unmarshal(gp, &gameplayer)
		if err != nil {
			return nil, err
		}
		gamePlayers = append(gamePlayers, gameplayer)
	}

	gameplayers := []vo.GamePlayer{}
	for _, p := range gamePlayers {

		player := vo.GamePlayer{
			Player: queries.Player{
				ID:        p.ID,
				Accountid: p.Accountid,
				Score:     p.Score,
				Isready:   p.Isready,
				Art:       p.Art,
			},
			Hand:  p.Hand,
			Play:  p.Play,
			Kitty: p.Kitty,
		}

		gameplayers = append(gameplayers, player)
	}

	match.Players = gameplayers

	return match, nil
}

func GetOpenMatches() ([]queries.Match, error) {
	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	matchesData, err := q.GetOpenMatches(ctx)

	if err != nil {
		return nil, err
	}

	var matches []queries.Match

	for _, matchData := range matchesData {
		var match queries.Match
		err = json.Unmarshal(matchData, &match)
		if err != nil {
			return nil, err
		}

		matches = append(matches, match)
	}

	return matches, nil
}

func CreateGameDeck(matchId *int) (*vo.GameDeck, error) {
	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//create a new deck object
	//create match cards based on queries.Cards
	//update deck_matchcard with new match card ids

	deck, err := q.GetDeckForMatchId(ctx, matchId)

	if err != nil {
		return &vo.GameDeck{}, err
	}

	cards, err := q.GetCards(ctx)

	if err != nil {
		return &vo.GameDeck{}, err
	}

	matchcards := []*vo.GameCard{}
	for _, card := range cards {
		matchCard, err := q.CreateMatchCards(ctx, queries.CreateMatchCardsParams{
			Cardid: card.ID,
			State:  queries.CardstateDeck,
		},
		)

		if err != nil {
			return &vo.GameDeck{}, err
		}

		gameCard := &vo.GameCard{
			Matchcard: matchCard,
			Card:      card,
		}

		matchcards = append(matchcards, gameCard)
	}

	gameDeck := &vo.GameDeck{
		Deck: &queries.Deck{
			ID: deck.ID,
		},
		Cards: matchcards,
	}

	return gameDeck, nil
}

func UpdateMatch(match queries.Match) error {
	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := q.UpdateMatch(ctx, queries.UpdateMatchParams{
		ID:                 match.ID,
		Creationdate:       match.Creationdate,
		Privatematch:       match.Privatematch,
		Elorangemin:        match.Elorangemin,
		Elorangemax:        match.Elorangemax,
		Deckid:             match.Deckid,
		Cutgamecardid:      match.Cutgamecardid,
		Currentplayerturn:  match.Currentplayerturn,
		Turnpasstimestamps: match.Turnpasstimestamps,
		Gamestate:          match.Gamestate,
		Art:                match.Art,
	})

	if err != nil {
		return err
	}

	return nil
}

func Deal() error {

	return nil
}
