package services

import (
	"encoding/json"
	"net/http"
	"strconv"

	"ca.openbracket.cribbage_cli/model"
	tea "github.com/charmbracelet/bubbletea"
)

func PostGame(a model.GameAction) tea.Msg {
	b, err := json.Marshal(a)
	if err != nil {
		return err
	}

	return url(EndPointGame, http.MethodPost, string(b))
}

func PollForReady(matchId int) tea.Msg {
	return url(EndPointGame, http.MethodGet, strconv.Itoa(matchId))
}
