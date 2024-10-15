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

func GetPlayerId(accountId int, players []model.Player) (*model.Player, error) {
	for _, p := range players {
		if p.AccountId == accountId {
			return &p, nil
		}
	}
	return nil, errors.New("no player found")
}

func CreateGame(accountId int) tea.Msg {
	newMatch := services.PostPlayerMatch(accountId).([]byte)

	var match *model.MatchDetailsResponse
	json.Unmarshal(newMatch, &match)
	return match
}

func GetPlayerForAccountId(id int, match *model.GameMatch) *model.Player {
	for _, player := range match.Players {
		if player.AccountId == id {
			return &player
		}
	}

	return nil
}
