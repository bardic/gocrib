package utils

import (
	"encoding/json"
	"errors"
	"slices"

	"github.com/bardic/gocrib/cli/services"
	"github.com/bardic/gocrib/model"
	"github.com/bardic/gocrib/queries"
	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/zap"
)

var Logger *zap.Logger

func GetCardById(id int32, deck *model.GameDeck) *queries.Card {
	for _, card := range deck.Cards {
		if card.Cardid == id {
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

func GetPlayerId(accountId int, players []queries.Player) (*queries.Player, error) {
	for _, p := range players {
		if int(p.Accountid) == accountId {
			return &p, nil
		}
	}
	return nil, errors.New("no player found")
}

func CreateGame(accountId int) tea.Msg {
	newMatch := services.PostPlayerMatch(accountId).([]byte)

	var matchDetails model.MatchDetailsResponse
	json.Unmarshal(newMatch, &matchDetails)

	return matchDetails
}

func GetPlayerForAccountId(id int, match *model.GameMatch) *queries.Player {
	for _, player := range match.Players {
		if int(player.Accountid) == id {
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
