package views

import (
	"github.com/bardic/cribbagev2/model"
)

type ViewModel struct {
	ViewState       model.ViewState
	Tabs            []string
	ActiveTab       int
	ActiveSlot      model.CardSlots
	Hand            []model.Card
	Kitty           []model.Card
	SelectedCardId  int
	SelectedCardIds []int
	CardsInPlay     []model.Card
}
