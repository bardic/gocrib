package utils

import (
	"encoding/json"
	"errors"
	"slices"
	"strings"

	"github.com/bardic/gocrib/cli/services"
	"github.com/bardic/gocrib/model"
	"github.com/bardic/gocrib/queries"
	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/zap"
)

var Logger *zap.Logger

func GetCardById(id int32, deck *model.GameDeck) *queries.Card {
	for _, card := range deck.Cards {
		if card.Matchcard.Cardid == id {
			return &card.Card
		}
	}
	return nil
}

func GetIdsFromCards(c []queries.Card) []int32 {
	ids := []int32{}
	for _, card := range c {
		ids = append(ids, card.ID)
	}

	return ids
}

func GetCardInHandById(id int32, hand []queries.Card) queries.Card {
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

func GetPlayerById(accountId int32, players []queries.Player) (*queries.Player, error) {
	for _, p := range players {
		if p.Accountid == accountId {
			return &p, nil
		}
	}
	return nil, errors.New("no player found")
}

func CreateGame(accountId int32) tea.Msg {
	newMatch := services.PostPlayerMatch(accountId).([]byte)

	var matchDetails model.MatchDetailsResponse
	json.Unmarshal(newMatch, &matchDetails)

	return matchDetails
}

func GetPlayerForAccountId(id int32, match *model.GameMatch) *queries.Player {
	for _, player := range match.Players {
		if player.Accountid == id {
			return &player
		}
	}

	return nil
}

func GetVisibleCards(activeTab int, player queries.Player) []int32 {
	var cards []int32
	switch activeTab {
	case 0:
		cards = nil
	case 1:
		cards = player.Play
	case 2:
		cards = player.Hand
	case 3:
		cards = player.Kitty
	}

	return cards
}

// func BuildPlayerInfo(player *queries.Player) string {
// 	return lipgloss.PlaceHorizontal(100, lipgloss.Right, lipgloss.NewStyle().PaddingRight(10).Render("\nPlayer Info\n"))
// }

func BuildFooter() string {
	return "\n\ntab/shift+tab: navigate screens • space: select • enter: submit • q: exit\n"
}

func IsPlayerTurn(playerId, matchId int32) bool {
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

func DrawRow(players []queries.Player, pegsToDraw, scoreOffet int) string {
	viewBuilder := strings.Builder{}
	for _, player := range players {
		viewBuilder.WriteString("\n")
		for i := range pegsToDraw {
			if i+scoreOffet == int(player.Score) {
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
