package utils

import (
	"slices"

	"github.com/bardic/gocrib/cli/state"
	"github.com/bardic/gocrib/model"
	"go.uber.org/zap"
)

var Logger *zap.Logger

func GetCardById(id int) *model.Card {
	for _, card := range state.ActiveDeck.Cards {
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
