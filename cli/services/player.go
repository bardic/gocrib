package services

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bardic/gocrib/cli/utils"
	"github.com/bardic/gocrib/vo"
	tea "github.com/charmbracelet/bubbletea"
)

type PReady struct {
	MatchId  *int
	PlayerId *int
}

func GetPlayer(playerId int) tea.Msg {
	id := int(playerId)
	pid := strconv.Itoa(id)

	return url(EndPointPlayer+"/?id="+pid, http.MethodGet, "")
}

func GetPlayerByForMatchAndAccount(matchId, playerId *int) tea.Msg {
	endpoint := utils.EndPointBuilder(EndPointPlayerByForMatchAndAccount, strconv.Itoa(*matchId), strconv.Itoa(*playerId))

	return url(endpoint, http.MethodGet, "")
}

func PutKitty(matchId, playerId *int, cards vo.HandModifier) tea.Msg {
	b, err := json.Marshal(cards)
	if err != nil {
		return err
	}

	endpoint := utils.EndPointBuilder(EndPointKitty, strconv.Itoa(*matchId), strconv.Itoa(*playerId))

	return url(endpoint, http.MethodPut, string(b))
}

func PutPlay(matchId, playerId *int, cards vo.HandModifier) tea.Msg {
	b, err := json.Marshal(cards)
	if err != nil {
		return err
	}

	endpoint := utils.EndPointBuilder(EndPointPlay, strconv.Itoa(*matchId), strconv.Itoa(*playerId))

	return url(endpoint, http.MethodPut, string(b))
}

func PostPlayer(accountId int) tea.Msg {
	return url(EndPointPlayer, http.MethodPost, strconv.Itoa(int(accountId)))
}

func DeletePlayer(ids []int) tea.Msg {
	return url(EndPointPlayer, http.MethodDelete, "")
}

func PlayerReady(playerId, matchId *int) tea.Msg {
	req := PReady{
		MatchId:  matchId,
		PlayerId: playerId,
	}

	reqStr, err := json.Marshal(req)

	if err != nil {
		return err
	}

	return url(EndPointPlayerReady, http.MethodPut, string(reqStr))
}
