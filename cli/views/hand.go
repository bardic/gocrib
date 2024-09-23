package views

import (
	"fmt"
	"slices"

	"github.com/bardic/cribbagev2/cli/styles"
	"github.com/bardic/cribbagev2/model"
	"github.com/charmbracelet/lipgloss"
)

func HandView(selectedCardId int, selectedCardIds []int, cards []model.Card) string {
	var s string

	cardViews := make([]string, 0)
	for i := 0; i < len(cards); i++ {
		view := fmt.Sprintf("%v : %v", cards[i].Suit, cards[i].Value)

		if slices.Contains(selectedCardIds, cards[i].Id) {
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
