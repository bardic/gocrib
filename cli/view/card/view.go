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

		top := lipgloss.Place(0, 0, lipgloss.Right, lipgloss.Bottom, cardStr)

		bottom := lipgloss.Place(8, 5, lipgloss.Right, lipgloss.Bottom, cardStr)

		styledCard = styles.ModelStyle.Render(top, bottom)

		if i == view.ActiveCardId {
			styledCard = styles.SelectedFocusedStyle.Render(top, bottom)
		}

		for _, v := range view.SelectedCardIds {
			if v == hand[i] {
				styledCard = styles.FocusedModelStyle.Render(top, bottom)
			}
		}

		//Original owner seems to be account id and not match player id

		if c.Match.Origowner != nil {
			styledCard = styles.PlayerStyles[*c.Match.Origowner].Render(styledCard)
		} else {
			styledCard = styles.PlayerStyles[*c.Match.Currowner].Render(styledCard)

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
