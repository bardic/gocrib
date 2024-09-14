package services

import (
	"encoding/json"
	"net/http"
	"strconv"

	"ca.openbracket.cribbage_cli/model"
	tea "github.com/charmbracelet/bubbletea"
)

func PostDeck(deck model.GameDeck) tea.Msg {
	b, err := json.Marshal(deck)
	if err != nil {
		return err
	}

	return url(EndPointDeck, http.MethodGet, string(b))
}

func PutDeck(deck model.GameDeck) tea.Msg {
	b, err := json.Marshal(deck)
	if err != nil {
		return err
	}

	return url(EndPointDeck, http.MethodPut, string(b))
}

func DeleteDeck(id int) tea.Msg {
	return url(EndPointDeckById+strconv.Itoa(id), http.MethodDelete, "")
}

func GetDeckById(id int) tea.Msg {
	return url(EndPointDeckById+strconv.Itoa(id), http.MethodGet, "")
}
