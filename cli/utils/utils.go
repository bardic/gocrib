package utils

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/bardic/gocrib/queries/queries"

	"github.com/bardic/gocrib/cli/styles"
	cliVO "github.com/bardic/gocrib/cli/vo"
	"github.com/bardic/gocrib/vo"

	"github.com/charmbracelet/lipgloss"
	"go.uber.org/zap"
)

var Logger *zap.Logger

func GetCardById(id int, deck *vo.GameDeck) *vo.GameCard {
	for _, card := range deck.Cards {
		gameCard := &vo.GameCard{
			Matchcard: card.Matchcard,
			Card:      card.Card,
		}

		if *card.Matchcard.Cardid == id {
			return gameCard
		}
	}
	return nil
}

func GetIdsFromCards(c []queries.Card) []*int {
	ids := []*int{}
	for _, card := range c {
		ids = append(ids, card.ID)
	}

	return ids
}

func GetCardInHandById(id *int, hand []queries.Card) queries.Card {
	idx := slices.IndexFunc(hand, func(c queries.Card) bool {
		return c.ID == id
	})

	if idx == -1 {
		return queries.Card{}
	}

	return hand[idx]
}

func NewLogger() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{
		"crib.log",
	}
	return cfg.Build()
}

func GetPlayerByAccountId(accountId *int, players []*queries.Player) (*queries.Player, error) {
	for _, p := range players {
		if p.Accountid == accountId {
			return p, nil
		}
	}
	return nil, errors.New("no player found")
}

// func CreateGame(accountId *int) vo.MatchDetailsResponse {
// 	newMatch := services.PostPlayerMatch(accountId).([]byte)

// 	var matchDetails vo.MatchDetailsResponse
// 	json.Unmarshal(newMatch, &matchDetails)

// 	return matchDetails
// }

func GetPlayerForAccountId(id *int, match *vo.GameMatch) *vo.GamePlayer {
	for _, player := range match.Players {
		if player.Accountid == id {
			return &player
		}
	}

	return nil
}

// func GetVisibleCards(activeTab int, player queries.Player) []*int {
// 	var cards []*int
// 	switch activeTab {
// 	case 0:
// 		cards = nil
// 	case 1:
// 		cards = player.Play
// 	case 2:
// 		cards = player.Hand
// 	case 3:
// 		cards = player.Kitty
// 	}

// 	return cards
// }

func BuildCommonFooter(activePlayerId, localPlayerId, matchId *int, gameState queries.Gamestate) string {
	f := fmt.Sprintf("\nState: %v | Local/Active Player : %v/%v | Match ID: %v ", gameState, localPlayerId, activePlayerId, matchId)
	f += "\ntab/shift+tab: navigate screens • space: select • enter: submit • q: exit\n"
	return f
}

func IsPlayerTurn(playerId, matchId int) bool {
	return playerId == matchId
}

func GetCardSuit(card *queries.Card) string {
	switch card.Suit {
	case queries.CardsuitSpades:
		return "♠"
	case queries.CardsuitHearts:
		return "♥"
	case queries.CardsuitDiamonds:
		return "♦"
	case queries.CardsuitClubs:
		return "♣"
	default:
		return "?"
	}
}

type PegState string

const (
	Filled PegState = "•"
	Empty  PegState = "○"
)

func DrawRow(players []vo.GamePlayer, pegsToDraw, scoreOffet int) string {
	viewBuilder := strings.Builder{}
	for _, player := range players {
		viewBuilder.WriteString("\n")

		for i := range pegsToDraw {
			if i+scoreOffet == *player.Score {
				viewBuilder.WriteString(string(Filled))
			} else {
				viewBuilder.WriteString(string(Empty))
			}
		}
	}

	viewBuilder.WriteString("\n")
	for range pegsToDraw {
		viewBuilder.WriteString("-")
	}

	return viewBuilder.String()
}

func RenderTabs(tabs []cliVO.Tab, activeTab int) []string {
	var renderedTabs []string

	for i, t := range tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(tabs)-1, i == activeTab
		if isActive {
			style = styles.ActiveTabStyle
		} else {
			style = styles.InactiveTabStyle
		}
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "└"
		} else if isLast && !isActive {
			border.BottomRight = "┴"
		}

		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(t.Title))
	}

	return renderedTabs
}

func GetPlayerIds(players []vo.GamePlayer) []*int {
	var playIds []*int
	for _, p := range players {
		playIds = append(playIds, p.ID)
	}
	return playIds
}

func IdFromCards(cards []queries.Matchcard) []int {
	var ids []int
	for _, c := range cards {
		ids = append(ids, *c.Cardid)
	}

	return ids
}

func EndPointBuilder(endpoint string, args ...string) string {
	for _, arg := range args {
		endpoint = strings.Replace(endpoint, "%s", arg, 1)
	}
	return endpoint
}
