package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/bardic/gocrib/cli/utils"

	tea "github.com/charmbracelet/bubbletea"
)

func GetPlayerMatchState(matchId *int) tea.Msg {
	id := strconv.Itoa(*matchId)
	return url(EndPointMatchState+"/?id="+id, http.MethodGet, "")
}

func GetMatchById(matchId *int) tea.Msg {
	id := strconv.Itoa(*matchId)
	return url(EndPointMatch+"/"+id, http.MethodGet, "")
}

func GetOpenMatches() tea.Msg {
	return url(EndPointOpenMatch, http.MethodGet, "")
}

func JoinMatch(accountId, activeMatchId int) tea.Msg {
	u := fmt.Sprintf(EndPointJoinMatch, activeMatchId, accountId)

	return url(u, http.MethodPut, "")
}

func PostPlayerMatch(accountId *int) tea.Msg {
	return sendReq(EndPointMatch+"/"+strconv.Itoa(*accountId), http.MethodPost, nil)
}

func CutDeck(matchId int, cutIndex string) tea.Msg {
	endpoint := utils.EndPointBuilder(EndPointMatchCutDeck, strconv.Itoa(matchId), cutIndex)

	return sendReq(endpoint, http.MethodPut, "")
}

func sendReq(endPoint string, method string, body interface{}) tea.Msg {
	b, err := json.Marshal(body)

	if err != nil {
		return err
	}

	return url(endPoint, method, string(b))
}
