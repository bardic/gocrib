package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/bardic/gocrib/vo"

	tea "github.com/charmbracelet/bubbletea"
)

func GetPlayerMatchState(matchId int32) tea.Msg {
	id := strconv.Itoa(int(matchId))
	return url(EndPointMatchState+"/?id="+id, http.MethodGet, "")
}

func GetPlayerMatch(matchId int32) tea.Msg {
	id := strconv.Itoa(int(matchId))
	return url(EndPointMatch+"/?id="+id, http.MethodGet, "")
}

func GetOpenMatches() tea.Msg {
	return url(EndPointOpenMatch, http.MethodGet, "")
}

func JoinMatch(accountId, activeMatchId int32) tea.Msg {
	u := fmt.Sprintf(EndPointJoinMatch, activeMatchId, accountId)

	return url(u, http.MethodPut, "")
}

func PostPlayerMatch(accountId int32) tea.Msg {
	req := vo.MatchRequirements{
		AccountId:   accountId,
		IsPrivate:   false,
		EloRangeMin: 1,
		EloRangeMax: 3000,
	}

	return sendReq(EndPointMatch, http.MethodPost, req)
}

func CutDeck(matchId int32, cutIndex string) tea.Msg {
	req := vo.CutDeckReq{
		MatchId:  matchId,
		CutIndex: cutIndex,
	}

	return sendReq(EndPointMatchCutDeck, http.MethodPut, req)
}

func sendReq(endPoint string, method string, body interface{}) tea.Msg {
	b, err := json.Marshal(body)

	if err != nil {
		return err
	}

	return url(endPoint, method, string(b))
}
