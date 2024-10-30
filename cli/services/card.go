package services

import (
	"net/http"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

func GetAllCards() tea.Msg {
	return url(EndPointAllCards, http.MethodGet, "")
}

func GetCardById(id int) tea.Msg {
	return url(EndPointCardById+strconv.Itoa(id), http.MethodGet, "")
}

func GetCardsForMatchId(matchId int32) tea.Msg {
	return url(EndPointGameplayCardsByMatch+strconv.Itoa(int(matchId)), http.MethodGet, "")
}
