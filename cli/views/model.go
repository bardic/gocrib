package views

import (
	"github.com/bardic/cribbagev2/model"
)

type ViewModel struct {
	ViewState      model.ViewState
	Tabs           []string
	ActiveTab      int
	ActiveSlot     model.CardSlots
	HighlighedId   int
	HighlightedIds []int
	CardsToDisplay []model.Card
}
