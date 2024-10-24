package kitty

import (
	"fmt"
	"slices"

	"github.com/bardic/gocrib/cli/styles"
	"github.com/bardic/gocrib/cli/utils"
	"github.com/bardic/gocrib/cli/views"
	"github.com/charmbracelet/lipgloss"
)

type KittyView struct {
	views.HandModel
}

func (v *KittyView) View() string {
	s := ""

	cardViews := make([]string, 0)
	for i := 0; i < len(v.Player.Kitty); i++ {
		c := utils.GetCardById(v.Player.Kitty[i], v.Deck)
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

func (v *KittyView) BuildHeader() string {
	return ""
}

func (v *KittyView) BuildFooter() string {
	return ""
}
