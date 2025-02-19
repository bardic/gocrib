package card

import (
	"fmt"

	"github.com/bardic/gocrib/cli/styles"
	"github.com/bardic/gocrib/cli/utils"
	"github.com/bardic/gocrib/vo"

	"github.com/charmbracelet/lipgloss"
)

type View struct {
	ActiveCardId    int
	SelectedCardIds []int
	// Match           *vo.GameMatch
	LocalPlayer *vo.GamePlayer
	// Deck            *vo.GameDeck
	Tabname string
	*vo.UIFooterVO
}

func NewCardView(match *vo.GameMatch, localPlayer *vo.GamePlayer, deck *vo.GameDeck, tabName string) *View {
	return &View{
		// Match:       match,
		LocalPlayer: localPlayer,
		// Deck:        deck,
		Tabname: tabName,
	}
}

func (view *View) Init() {
	// view.ActiveCardId = 0
}

func (view *View) Update(gameMatch *vo.GameMatch) error {
	return nil
}

func (view *View) Render(gameMatch *vo.GameMatch, gameDeck *vo.GameDeck, hand []int) string {
	var s string
	var cardViews []string

	s += view.BuildHeader()
	for i := 0; i < len(hand); i++ {
		c := utils.GetCardById(hand[i], gameDeck)

		if c == nil {
			continue
		}

		cardStr := fmt.Sprintf("%v%v", utils.GetCardSuit(&c.Card), c.Card.Value)
		var styledCard string

		styledCard = styles.ModelStyle.Render(cardStr)

		if i == view.ActiveCardId {
			styledCard = styles.SelectedFocusedStyle.Render(cardStr)
		}

		for _, v := range view.SelectedCardIds {
			if v == hand[i] {
				styledCard = styles.FocusedModelStyle.Render(cardStr)
			}
		}

		cardViews = append(cardViews, styledCard)
	}

	s += lipgloss.JoinHorizontal(lipgloss.Top, cardViews...)

	s += "\n" + view.BuildFooter(gameMatch)

	return s
}

func (view *View) BuildHeader() string {
	return ""
}

func (view *View) BuildFooter(gameMatch *vo.GameMatch) string {
	f := utils.BuildCommonFooter(gameMatch, view.LocalPlayer)

	return f
}
