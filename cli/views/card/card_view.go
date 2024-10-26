package card

import (
	"fmt"
	"slices"

	"cli/styles"
	"cli/utils"
	"cli/views"

	"github.com/charmbracelet/lipgloss"
)

type CardView struct {
	*views.HandModel
	SelectedCardId int
}

func (v *CardView) Init() {
	v.SelectedCardId = 0
}

func (v *CardView) Render() string {
	s := ""
	cardViews := make([]string, 0)

	for i := 0; i < len(v.CardsToDisplay); i++ {
		c := utils.GetCardById(v.CardsToDisplay[i], v.Deck)
		view := fmt.Sprintf("%v%v", utils.GetCardSuit(c), c.Value)

		if slices.Contains(v.SelectedCardIds, c.ID) {
			if i == v.SelectedCardId {
				cardViews = append(cardViews, styles.SelectedFocusedStyle.Render(view))
			} else {
				cardViews = append(cardViews, styles.SelectedStyle.Render(view))
			}
		} else {
			if i == v.SelectedCardId {
				cardViews = append(cardViews, styles.FocusedModelStyle.Render(view))
			} else {
				cardViews = append(cardViews, styles.ModelStyle.Render(view))
			}
		}
	}

	s += lipgloss.JoinHorizontal(lipgloss.Top, cardViews...)
	return s
}

func (v *CardView) BuildHeader() string {
	return ""
}

func (v *CardView) BuildFooter() string {
	return ""
}
