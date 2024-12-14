package services

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bardic/gocrib/vo"
	tea "github.com/charmbracelet/bubbletea"
)

type PReady struct {
	MatchId  int32
	PlayerId int32
}

func GetPlayer(playerId int32) tea.Msg {
	id := int(playerId)
	pid := strconv.Itoa(id)

	return url(EndPointPlayer+"/?id="+pid, http.MethodGet, "")
}

func PutKitty(req vo.HandModifier) tea.Msg {
	b, err := json.Marshal(req)
	if err != nil {
		return err
	}

	return url(EndPointKitty, http.MethodPut, string(b))
}

func PutPlay(req vo.HandModifier) tea.Msg {
	b, err := json.Marshal(req)
	if err != nil {
		return err
	}

	return url(EndPointPlay, http.MethodPut, string(b))
}

func PostPlayer(accountId int32) tea.Msg {
	return url(EndPointPlayer, http.MethodPost, strconv.Itoa(int(accountId)))
}

func DeletePlayer(ids []int) tea.Msg {
	return url(EndPointPlayer, http.MethodDelete, "")
}

func PlayerReady(playerId int32, matchId int32) tea.Msg {
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
