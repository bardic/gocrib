package services

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func GetAllCards() tea.Msg {
	return url(EndPointAllCards, http.MethodGet, "")
}

func GetCardById(id int) tea.Msg {
	return url(EndPointCardById+strconv.Itoa(id), http.MethodGet, "")
}

func GetGameplayCardsByIds(ids []int) tea.Msg {
	s, _ := json.Marshal(ids)
	return url(EndPointGameplayCardsByIds+strings.Trim(string(s), "[]"), http.MethodGet, "")
}

func GetGampleCardsForMatch(matchId int32) tea.Msg {
	return url(EndPointGameplayCardsByMatch+strconv.Itoa(int(matchId)), http.MethodGet, "")
}
