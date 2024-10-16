package utils

import (
	"encoding/json"
	"errors"
	"slices"

	"github.com/bardic/gocrib/cli/services"
	"github.com/bardic/gocrib/model"
	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/zap"
)

var Logger *zap.Logger

func GetCardById(id int, deck *model.GameDeck) *model.Card {
	for _, card := range deck.Cards {
		if card.CardId == id {
			return &card.Card
		}
	}
	return nil
}

func GetIdsFromCards(c []model.Card) []int {
	ids := []int{}
	for _, card := range c {
		ids = append(ids, card.Id)
	}

	return ids
}

func GetCardInHandById(id int, hand []model.Card) model.Card {
	idx := slices.IndexFunc(hand, func(c model.Card) bool {
		return c.Id == id
	})

	if idx == -1 {
		return model.Card{}
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
		if p.AccountId == accountId {
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

func GetPlayerForAccountId(id int, match *queries.Match) *queries.Player {
	for _, player := range match.Players {
		if player.AccountId == id {
			return &player
		}
	}

	return nil
}

func GetVisibleCards(activeTab int, player queries.Player) []int {
	var cards []int
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
