package utils

import (
	"fmt"
	"slices"
	"strings"

	"github.com/bardic/gocrib/queries/queries"

	"github.com/bardic/gocrib/vo"

	"go.uber.org/zap"
)

var Logger *zap.Logger

func GetCardById(id int, deck *vo.GameDeck) *vo.GameCard {
	for _, card := range deck.Cards {
		gameCard := &vo.GameCard{
			Match: card.Match,
			Card:  card.Card,
		}

		if *card.Match.Cardid == id {
			return gameCard
		}
	}
	return nil
}

func GetIdsFromCards(c []queries.Card) []*int {
	ids := []*int{}
	for _, card := range c {
		ids = append(ids, card.ID)
	}

	return ids
}

func GetCardInHandById(id *int, hand []queries.Card) queries.Card {
	idx := slices.IndexFunc(hand, func(c queries.Card) bool {
		return c.ID == id
	})

	if idx == -1 {
		return queries.Card{}
	}

	return hand[idx]
}

func NewLogger() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{
		"crib.log",
	}
	return cfg.Build()
}

func GetPlayerForAccountId(id *int, match *vo.GameMatch) *vo.GamePlayer {
	for _, player := range match.Players {
		if *player.Accountid == *id {
			return player
		}
	}

	return nil
}

func BuildCommonFooter(match *vo.GameMatch, localplayer *vo.GamePlayer) string {
	localPlayerId := "-"
	activePlayerId := "-"

	if localplayer != nil {
		localPlayerId = fmt.Sprintf("%v", &localplayer.ID)
		activePlayerId = fmt.Sprintf("%v", &match.Currentplayerturn)
	}

	f := fmt.Sprintf("\nState: %v | Local/Active Player : %v/%v | Match ID: %v ", match.Gamestate, localPlayerId, activePlayerId, fmt.Sprintf("%v", &match.ID))
	f += "\ntab/shift+tab: navigate screens • space: select • enter: submit • q: exit\n"
	return f
}

func IsPlayerTurn(playerId, matchId int) bool {
	return playerId == matchId
}

func GetCardSuit(card *queries.Card) string {
	switch card.Suit {
	case queries.CardsuitSpades:
		return "♠"
	case queries.CardsuitHearts:
		return "♥"
	case queries.CardsuitDiamonds:
		return "♦"
	case queries.CardsuitClubs:
		return "♣"
	default:
		return "?"
	}
}

type PegState string

const (
	Filled PegState = "•"
	Empty  PegState = "○"
)

func DrawRow(players []*vo.GamePlayer, pegsToDraw, scoreOffet int) string {
	viewBuilder := strings.Builder{}
	for _, player := range players {
		viewBuilder.WriteString("\n")

		for i := range pegsToDraw {
			if i+scoreOffet == *player.Score {
				viewBuilder.WriteString(string(Filled))
			} else {
				viewBuilder.WriteString(string(Empty))
			}
		}
	}

	viewBuilder.WriteString("\n")
	for range pegsToDraw {
		viewBuilder.WriteString("-")
	}

	return viewBuilder.String()
}

func GetPlayerIds(players []*vo.GamePlayer) []*int {
	var playIds []*int
	for _, p := range players {
		playIds = append(playIds, p.ID)
	}
	return playIds
}

func IdFromCards(cards []queries.Matchcard) []int {
	var ids []int
	for _, c := range cards {
		ids = append(ids, *c.Cardid)
	}

	return ids
}

func EndPointBuilder(endpoint string, args ...string) string {
	for _, arg := range args {
		endpoint = strings.Replace(endpoint, "%s", arg, 1)
	}
	return endpoint
}
