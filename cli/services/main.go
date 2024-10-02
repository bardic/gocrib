package services

import (
	"bytes"
	"io"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	BaseUrl = "http://localhost:1323/v1"
	//Cards
	EndPointAllCards           = BaseUrl + "/player/allcards"
	EndPointCardById           = BaseUrl + "/player/card?id="
	EndPointGameplayCardsByIds = BaseUrl + "/player/gameplaycards?ids="
	//Game
	EndPointGame = BaseUrl + "/game/playCard"
	//Player
	EndPointPlayer = BaseUrl + "/player/player"
	EndPointKitty  = BaseUrl + "/player/kitty"
	EndPointPlay   = BaseUrl + "/player/play"
	//History
	EndPointHistory = BaseUrl + "/history"
	//Match
	EndPointMatch     = BaseUrl + "/player/match"
	EndPointJoinMatch = BaseUrl + "/player/match/join"
	EndPointMatches   = BaseUrl + "/player/matches"
	EndPointOpenMatch = BaseUrl + "/player/matches/open"
	EndPointMatchCard = BaseUrl + "/player/match/card"
	//Chat
	EndPointChat = BaseUrl + "/chat"
	//Lobby
	EndPointLobby = BaseUrl + "/player/lobby"
	//Deck
	EndPointDeck     = BaseUrl + "/deck"
	EndPointDeckById = BaseUrl + "/player/match/deck?id="
	//Account
	EndPointLogin = BaseUrl + "/account/login"
)

func url(url string, method string, json string) tea.Msg {
	var buf *bytes.Buffer
	buf = bytes.NewBuffer([]byte(json))

	if json == "" {
		buf = bytes.NewBuffer(nil)
	}

	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	return body

}
