package services

import (
	"encoding/json"
	"net/http"
	"strconv"

	"vo"

	tea "github.com/charmbracelet/bubbletea"
)

// func GetMatchesForPlayerId() tea.Msg {
// 	return url(EndPointMatches+"/?id="+strconv.Itoa(state.AccountId), http.MethodGet, "")
// }

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

func JoinMatch(playerId, activeMatchId int) tea.Msg {
	req := vo.JoinMatchReq{
		PlayerId: playerId,
		MatchId:  activeMatchId,
	}

	b, err := json.Marshal(req)

	if err != nil {
		return err
	}

	return url(EndPointJoinMatch, http.MethodPut, string(b))
}

func PutPlayerMatch(id int) tea.Msg {
	return url(EndPointMatch, http.MethodPut, "")
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
