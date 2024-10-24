package hand

import (
	"fmt"
	"slices"

	"github.com/bardic/gocrib/cli/styles"
	"github.com/bardic/gocrib/cli/utils"
	"github.com/bardic/gocrib/cli/views"
	"github.com/charmbracelet/lipgloss"
)

type HandView struct {
	views.HandModel
}

func (v *HandView) View() string {
	s := ""

	cardViews := make([]string, 0)
	for i := 0; i < len(v.Player.Hand); i++ {
		c := utils.GetCardById(v.Player.Hand[i], v.Deck)
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

func (v *HandView) BuildHeader() string {
	return ""
}

func (v *HandView) BuildFooter() string {
	return ""
}
