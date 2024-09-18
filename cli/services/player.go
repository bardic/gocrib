package services

import (
	"encoding/json"
	"net/http"

	"github.com/bardic/cribbagev2/model"
	tea "github.com/charmbracelet/bubbletea"
)

func GetPlayer() tea.Msg {
	return url(EndPointPlayer, http.MethodGet, "")
}

func PutPlayer(id int) tea.Msg {
	return url(EndPointPlayer, http.MethodPut, "")
}

func PostPlayer(player model.Player) tea.Msg {

	b, err := json.Marshal(player)
	if err != nil {
		return err
	}

	return url(EndPointPlayer, http.MethodPost, string(b))
}

func DeletePlayer(ids []int) tea.Msg {
	return url(EndPointPlayer, http.MethodDelete, "")
}
