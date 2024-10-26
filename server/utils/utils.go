package utils

import (
	"context"
	"encoding/json"
	"math/rand/v2"

	"queries"
	conn "server/db"

	"vo"
)

func QueryForCards(ids []int32) ([]vo.GameCard, error) {
	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx := context.Background()

	matchCards, err := q.GetMatchCards(ctx, ids)

	if err != nil {
		return []vo.GameCard{}, err
	}

	baseCards, err := q.GetCards(ctx)

	if err != nil {
		return []vo.GameCard{}, err
	}

	cards := []vo.GameCard{}

	for _, matchCard := range matchCards {
		card := GetCardByIdFromCards(int(matchCard.ID), baseCards)
		cards = append(cards, vo.GameCard{
			Matchcard: matchCard,
			Card:      card,
		})
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

func UpdatePlay(details vo.HandModifier) (*vo.GameMatch, error) {
	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx := context.Background()

	q.UpdateCardsPlayed(ctx, queries.UpdateCardsPlayedParams{
		Play: details.CardIds,
		ID:   details.PlayerId,
	})

	return PlayCard(details)
}

func PlayCard(details vo.HandModifier) (*vo.GameMatch, error) {
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

	var match vo.GameMatch
	err = json.Unmarshal(b, &match)
	if err != nil {
		return nil, err
	}

	return &match, nil
}

func PlayersReady(players []*queries.Player) bool {
	ready := true

	// if len(players) < 2 {
	// 	return false
	// }

	// for _, p := range players {
	// 	if !p.Isready {
	// 		ready = false
	// 	}
	// }

	return ready
}

func GetMatchForPlayerId(playerId int) (*vo.GameMatch, error) {
	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx := context.Background()

	b, err := q.GetMatchByPlayerId(ctx, int32(playerId))

	if err != nil {
		return nil, err
	}

	var match *vo.GameMatch
	err = json.Unmarshal(b, &match)
	if err != nil {
		return nil, err
	}

	return match, nil
}

func GetDeckById(id int32) (queries.Deck, error) {
	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx := context.Background()

	d, err := q.GetDeck(ctx, id)

	if err != nil {
		return queries.Deck{}, err
	}

	return d, nil
}

func Deal(match *vo.GameMatch) (*queries.Deck, error) {
	deck, err := GetDeckById(match.Deckid)

	// deck = *deck.Shuffle()
	if err != nil {
		return nil, err
	}

	cardsPerHand := 6
	if len(match.Players) == 3 {
		cardsPerHand = 5
	}

	for i := 0; i < len(match.Players)*cardsPerHand; i++ {
		var cardId int32
		cardId, deck.Cards = deck.Cards[0], deck.Cards[1:]
		idx := len(match.Players) - 1 - i%len(match.Players)

		if len(match.Players[idx].Hand) < cardsPerHand {
			match.Players[idx].Hand = append(match.Players[idx].Hand, cardId)
		}
	}

	for _, p := range match.Players {
		_, err := UpdatePlayerById(p)

		if err != nil {
			return nil, err
		}
	}

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

func UpdatePlayerById(player *queries.Player) (queries.Player, error) {
	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx := context.Background()

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

	return true
}
