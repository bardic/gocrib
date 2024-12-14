package services

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bardic/gocrib/vo"

	tea "github.com/charmbracelet/bubbletea"
)

func PostGame() tea.Msg {
	a := vo.GameAction{}
	b, err := json.Marshal(a)
	if err != nil {
		return err
	}

	return url(EndPointGame, http.MethodPost, string(b))
}

func PollForReady(matchId int) tea.Msg {
	return url(EndPointGame, http.MethodGet, strconv.Itoa(matchId))
}
