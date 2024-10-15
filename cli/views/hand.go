package views

import (
	"fmt"
	"slices"

	"github.com/bardic/gocrib/cli/styles"
	"github.com/bardic/gocrib/cli/utils"
	"github.com/bardic/gocrib/model"
	"github.com/charmbracelet/lipgloss"
)

// function to get card suit
func getCardSuit(card *model.Card) string {
	switch card.Suit {
	case model.Spades:
		return "♠"
	case model.Hearts:
		return "♥"
	case model.Diamonds:
		return "♦"
	case model.Clubs:
		return "♣"
	default:
		return "?"
	}
}

func HandView(selectedCardId int, selectedCardIds []int, cards []int, deck *model.GameDeck) string {
	var s string

	cardViews := make([]string, 0)
	for i := 0; i < len(cards); i++ {
		c := utils.GetCardById(cards[i], deck)
		view := fmt.Sprintf("%v%v", getCardSuit(c), c.Value)

		if slices.Contains(selectedCardIds, c.Id) {
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
