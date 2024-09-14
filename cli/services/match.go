package services

import (
	"encoding/json"
	"net/http"

	"ca.openbracket.cribbage_cli/model"
	tea "github.com/charmbracelet/bubbletea"
)

func GetPlayerMatch() tea.Msg {
	return url(EndPointMatch, http.MethodGet, "")
}

func PutPlayerMatch(id int) tea.Msg {
	return url(EndPointMatch, http.MethodPut, "")
}

func PostPlayerMatch(match model.Match) tea.Msg {
	b, err := json.Marshal(match)
	if err != nil {
		return err
	}

	return url(EndPointMatch, http.MethodPost, string(b))
}

func GetPlayerMatchCard() tea.Msg {
	return url(EndPointMatchCard, http.MethodPut, "")
}

func PutPlayerMatchCard(id int) tea.Msg {
	return url(EndPointMatchCard, http.MethodPut, "")
}

func PostPlayerMatchCard(ids []int) tea.Msg {
	return url(EndPointMatchCard, http.MethodPost, "")
}
