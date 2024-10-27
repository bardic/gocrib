package card

import (
	"fmt"
	"slices"

	"cli/styles"
	"cli/utils"
	cliVO "cli/vo"

	"github.com/charmbracelet/lipgloss"
)

type View struct {
	*cliVO.HandVO
	ActiveCardId    int32
	SelectedCardIds []int32
}

func (view *View) Init() {
	view.ActiveCardId = 0
}

func (view *View) Render() string {
	var s string
	var cardViews []string

	for i := 0; i < len(view.CardIds); i++ {
		c := utils.GetCardById(view.CardIds[i], view.Deck)
		cardStr := fmt.Sprintf("%v%v", utils.GetCardSuit(c), c.Value)

		if slices.Contains(view.SelectedCardIds, c.ID) {
			if int32(i) == view.ActiveCardId {
				cardViews = append(cardViews, styles.SelectedFocusedStyle.Render(cardStr))
			} else {
				cardViews = append(cardViews, styles.SelectedStyle.Render(cardStr))
			}
		} else {
			if int32(i) == view.ActiveCardId {
				cardViews = append(cardViews, styles.FocusedModelStyle.Render(cardStr))
			} else {
				cardViews = append(cardViews, styles.ModelStyle.Render(cardStr))
			}
		}
	}

	s += lipgloss.JoinHorizontal(lipgloss.Top, cardViews...)
	return s
}

func (view *View) BuildHeader() string {
	return ""
}

func (view *View) BuildFooter() string {
	return ""
}
