package utils

import (
	"context"
	"encoding/json"
	"math/rand/v2"

	"github.com/bardic/gocrib/queries"
	conn "github.com/bardic/gocrib/server/db"

	"github.com/bardic/gocrib/model"
)

func QueryForCards(ids []int32) ([]model.GameCard, error) {
	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx := context.Background()

	matchCards, err := q.GetMatchCards(ctx, ids)
	baseCards, err := q.GetCards(ctx)

	cards := []model.GameCard{}

	for _, matchCard := range matchCards {
		card := GetCardByIdFromCards(int(matchCard.ID), baseCards)
		cards = append(cards, model.GameCard{
			Matchcard: matchCard,
			Card:      card,
		})
	}

	if err != nil {
		return []model.GameCard{}, err
	}

	return cards, nil
}

func GetCardByIdFromCards(cardId int, cards []queries.Card) queries.Card {
	for _, c := range cards {
		if c.ID == int32(cardId) {
			return c
		}
	}

	return queries.Card{}

}

func UpdatePlay(details model.HandModifier) (*model.GameMatch, error) {
	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx := context.Background()

	q.UpdateCardsPlayed(ctx, queries.UpdateCardsPlayedParams{
		Play: details.CardIds,
		ID:   details.MatchId,
	})

	return PlayCard(details)
}

func PlayCard(details model.HandModifier) (*model.GameMatch, error) {
	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx := context.Background()

	q.RemoveCardsFromHand(ctx, queries.RemoveCardsFromHandParams{
		ID:   details.PlayerId,
		Hand: details.CardIds,
	})

	b, err := q.GetMatchById(ctx, details.MatchId)

	if err != nil {
		return nil, err
	}

	var match model.GameMatch
	err = json.Unmarshal(b, &match)
	if err != nil {
		return nil, err
	}

	return &match, nil
}

func PlayersReady(players []queries.Player) bool {
	ready := true

	if len(players) < 2 {
		return false
	}

	for _, p := range players {
		if !p.Isready {
			ready = false
		}
	}

	return ready
}

func GetMatchForPlayerId(playerId int) (int, error) {
	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx := context.Background()

	b, err := q.GetMatchByPlayerId(ctx, int32(playerId))

	if err != nil {
		return 0, err
	}

	var match model.GameMatch
	err = json.Unmarshal(b, &match)
	if err != nil {
		return 0, err
	}

	return int(match.ID), nil
}

func Deal(match *queries.Match) (*queries.Deck, error) {
	deck, err := GetDeckById(int(match.Deckid))

	// deck = *deck.Shuffle()
	if err != nil {
		return nil, err
	}

	// players := []queries.Player{}

	// for _, id := range match.Playerids {
	// 	player, err := GetPlayerById(int(id))

	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	players = append(players, player)
	// }

	// cardsPerHand := 6
	// if len(players) == 3 {
	// 	cardsPerHand = 5
	// }

	// for i := 0; i < len(players)*cardsPerHand; i++ {
	// 	var card queries.Card
	// 	card, deck.Cards = deck.Cards[0], deck.Cards[1:]
	// 	idx := len(players) - 1 - i%len(players)

	// 	if len(players[idx].Hand) < cardsPerHand {
	// 		players[idx].Hand = append(players[idx].Hand, card.Cardid)
	// 	}
	// }

	// for _, p := range players {
	// 	UpdatePlayerById(p)
	// }

	return &deck, nil
}

func GetPlayerById(id int) (queries.Player, error) {
	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx := context.Background()

	p, err := q.GetPlayer(ctx, int32(id))

	if err != nil {
		return queries.Player{}, err
	}

	return p, nil
}

func UpdatePlayerById(player queries.Player) (queries.Player, error) {
	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx := context.Background()

	err := q.UpdatePlayer(ctx, queries.UpdatePlayerParams{
		Hand:    player.Hand,
		Play:    player.Play,
		Kitty:   player.Kitty,
		Score:   player.Score,
		Isready: player.Isready,
		Art:     player.Art,
	})

	if err != nil {
		return queries.Player{}, err
	}

	return player, nil
}

func Shuffle(d *queries.Deck) *queries.Deck {
	rand.Shuffle(len(d.Cards), func(i, j int) {
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	})

	return d
}

func Eq(p *queries.Player, c *queries.Player) bool {
	if p.ID != c.ID {
		return false
	}

	if p.Accountid != c.Accountid {
		return false
	}

	if p.Score != c.Score {
		return false
	}

	if p.Art != c.Art {
		return false
	}

	// if !eqIntArr(p.Hand, c.Hand) {
	// 	return false
	// }

	// if !eqIntArr(p.Play, c.Play) {
	// 	return false
	// }

	// if !eqIntArr(p.Kitty, c.Kitty) {
	// 	return false
	// }

	return true
}
