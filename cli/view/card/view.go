package card

import (
	"fmt"
	"slices"

	"github.com/bardic/gocrib/cli/styles"
	"github.com/bardic/gocrib/cli/utils"
	"github.com/bardic/gocrib/vo"

	"github.com/charmbracelet/lipgloss"
)

type View struct {
	ActiveCardId    int
	SelectedCardIds []int
	Match           *vo.GameMatch
	LocalPlayer     *vo.GamePlayer
	Deck            *vo.GameDeck
	*vo.UIFooterVO
}

func NewCardView(match *vo.GameMatch, localPlayer *vo.GamePlayer, deck *vo.GameDeck) *View {
	return &View{
		Match:       match,
		LocalPlayer: localPlayer,
		Deck:        deck,
	}
}

func (view *View) Init() {
	// view.ActiveCardId = 0
}

func (view *View) Render(hand []int) string {
	var s string
	var cardViews []string

	s += view.BuildHeader()
	for i := 0; i < len(hand); i++ {
		c := utils.GetCardById(hand[i], view.Deck)

		if c == nil {
			continue
		}

		cardStr := fmt.Sprintf("%v%v", utils.GetCardSuit(&c.Card), c.Card.Value)
		styledCard := styles.ModelStyle.Render(cardStr)
		if slices.Index(view.SelectedCardIds, *c.Match.ID) > -1 {
			if int(i) == view.ActiveCardId {
				styledCard = styles.SelectedFocusedStyle.Render(cardStr)
			} else {
				styledCard = styles.FocusedModelStyle.Render(cardStr)
			}
		} else {
			if int(i) == view.ActiveCardId {
				styledCard = styles.FocusedModelStyle.Render(cardStr)
			}
		}
		cardViews = append(cardViews, styledCard)
	}

	s += lipgloss.JoinHorizontal(lipgloss.Top, cardViews...)

	s += "\n" + view.BuildFooter()

	return s
}

func (view *View) BuildHeader() string {
	return ""
}

func (view *View) BuildFooter() string {
	f := utils.BuildCommonFooter(view.Match, view.LocalPlayer)

	return f
}
