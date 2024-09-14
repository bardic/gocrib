package services

import (
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
)

func GetChat() tea.Msg {
	return url(EndPointChat, http.MethodGet, "")
}

func PutChat(id int) tea.Msg {
	return url(EndPointChat, http.MethodPut, "")
}

func PostChat(ids []int) tea.Msg {
	return url(EndPointChat, http.MethodPost, "")
}

func DeleteChat(ids []int) tea.Msg {
	return url(EndPointChat, http.MethodDelete, "")

}
