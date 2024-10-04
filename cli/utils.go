package main

import (
	"slices"

	"github.com/bardic/cribbagev2/cli/state"
	"github.com/bardic/cribbagev2/model"
)

func getCardById(id int) *model.Card {
	for _, card := range state.ActiveDeck.Cards {
		if card.CardId == id {
			return &card.Card
		}
	}
	return nil
}

func getIdsFromCards(c []model.Card) []int {
	ids := make([]int, len(c))
	for _, card := range c {
		ids = append(ids, card.Id)
	}

	return ids
}

func getCardInHandById(id int, hand []model.Card) model.Card {
	idx := slices.IndexFunc(hand, func(c model.Card) bool {
		return c.Id == id
	})

	if idx == -1 {
		return model.Card{}
	}

	return hand[idx]
}
