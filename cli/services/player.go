package services

import (
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
)

func GetPlayer() tea.Msg {
	return url(EndPointPlayer, http.MethodGet, "")
}

func PutPlayer(id int) tea.Msg {
	return url(EndPointPlayer, http.MethodPut, "")
}

func PostPlayer(ids []int) tea.Msg {
	return url(EndPointPlayer, http.MethodPost, "")
}

func DeletePlayer(ids []int) tea.Msg {
	return url(EndPointPlayer, http.MethodDelete, "")
}
