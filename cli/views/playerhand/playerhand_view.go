package playerhand

import (
	"fmt"
	"slices"

	"cli/styles"
	"cli/utils"
	"cli/views"

	"github.com/charmbracelet/lipgloss"
)

type PlayerHandView struct {
	views.HandModel
	SelectedCardId int
}

func (v *PlayerHandView) Init() {
	v.SelectedCardId = 0
}

func (v *PlayerHandView) Render() string {
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

func (v *PlayerHandView) BuildHeader() string {
	return ""
}

func (v *PlayerHandView) BuildFooter() string {
	return ""
}
