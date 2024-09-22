package views

import (
	"github.com/bardic/cribbagev2/model"
)

type ViewModel struct {
	ActiveSlot      model.CardSlots
	GameState       model.GameState
	ViewState       model.ViewState
	Hand            []model.Card
	Kitty           []model.Card
	CardsInPlay     []model.Card
	SelectedCardIds []int
	Next            int
	ActiveTab       int
	Tabs            []string
}
