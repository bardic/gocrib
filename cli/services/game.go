package services

import (
	"encoding/json"
	"model"
	"net/http"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

func PostGame() tea.Msg {
	a := model.GameAction{}
	b, err := json.Marshal(a)
	if err != nil {
		return err
	}

	return url(EndPointGame, http.MethodPost, string(b))
}

func PollForReady(matchId int) tea.Msg {
	return url(EndPointGame, http.MethodGet, strconv.Itoa(matchId))
}
