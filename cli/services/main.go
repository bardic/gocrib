// Services are the CLIs connection to the backend. If you look in server/route you'll see a coorsponding package for each service.
package services

import (
	"bytes"
	"io"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	BaseUrl = "http://localhost:1323/v1"
	//Game
	EndPointGame = BaseUrl + "/game/playCard"
	//Player
	EndPointPlayer      = BaseUrl + "/match/player"
	EndPointKitty       = BaseUrl + "/match/player/kitty"
	EndPointPlay        = BaseUrl + "/match/player/play"
	EndPointPlayerReady = BaseUrl + "/match/player/ready"
	//Match
	EndPointMatch                = BaseUrl + "/match"
	EndPointMatchState           = BaseUrl + "/match/state"
	EndPointJoinMatch            = BaseUrl + "/match/join"
	EndPointOpenMatch            = BaseUrl + "/match/open"
	EndPointMatchCard            = BaseUrl + "/match/card"
	EndPointGameplayCardsByMatch = BaseUrl + "/match/cards?id="
	EndPointMatchCutDeck         = BaseUrl + "/match/cut"
	EndPointDeckById             = BaseUrl + "/match/deck?id="
	//Deck
	EndPointDeck = BaseUrl + "/deck"
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
