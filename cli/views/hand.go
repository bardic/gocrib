package views

import (
	"fmt"
	"slices"

	"github.com/bardic/gocrib/cli/styles"
	"github.com/bardic/gocrib/cli/utils"
	"github.com/bardic/gocrib/model"
	"github.com/bardic/gocrib/queries"
	"github.com/charmbracelet/lipgloss"
)

// function to get card suit
func getCardSuit(card *queries.Card) string {
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

func HandView(selectedCardId int, selectedCardIds []int32, cards []int32, deck *model.GameDeck) string {
	var s string

	cardViews := make([]string, 0)
	for i := 0; i < len(cards); i++ {
		c := utils.GetCardById(cards[i], deck)
		view := fmt.Sprintf("%v%v", getCardSuit(c), c.Value)

		if slices.Contains(selectedCardIds, c.ID) {
			if i == selectedCardId {
				cardViews = append(cardViews, styles.SelectedFocusedStyle.Render(view))
			} else {
				cardViews = append(cardViews, styles.SelectedStyle.Render(view))
			}
		} else {
			if i == selectedCardId {
				cardViews = append(cardViews, styles.FocusedModelStyle.Render(view))
			} else {
				cardViews = append(cardViews, styles.ModelStyle.Render(view))
			}
		}
	}

	s += lipgloss.JoinHorizontal(lipgloss.Top, cardViews...)
	s += styles.HelpStyle.Render("\ntab: focus next • q: exit\n")
	return s
}
