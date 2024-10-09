package services

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bardic/gocrib/cli/state"
	"github.com/bardic/gocrib/model"
	tea "github.com/charmbracelet/bubbletea"
)

func GetMatchesForPlayerId() tea.Msg {
	return url(EndPointMatches+"/?id="+strconv.Itoa(state.AccountId), http.MethodGet, "")
}

func GetPlayerMatch() tea.Msg {
	return url(EndPointMatch+"/?id="+strconv.Itoa(state.ActiveMatchId), http.MethodGet, "")
}

func GetOpenMatches() tea.Msg {
	return url(EndPointOpenMatch, http.MethodGet, "")
}

func JoinMatch() tea.Msg {
	req := model.JoinMatchReq{
		RequesterId: state.AccountId,
		MatchId:     state.ActiveMatchId,
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

func PostPlayerMatch() tea.Msg {
	req := model.MatchRequirements{
		RequesterId: state.AccountId,
		IsPrivate:   false,
		EloRangeMin: 1,
		EloRangeMax: 3000,
	}

	b, err := json.Marshal(req)

	if err != nil {
		return err
	}

	return url(EndPointMatch, http.MethodPost, string(b))
}
