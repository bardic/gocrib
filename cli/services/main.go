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
	EndPointPlayer                     = BaseUrl + "/match/%s/player"
	EndPointKitty                      = BaseUrl + "/match/%s/player/%s/to/%s/kitty"
	EndPointPlayerByForMatchAndAccount = BaseUrl + "/match/%s/account/%s"
	EndPointPlay                       = BaseUrl + "/match/%s/player/%s/to/%s/play"
	EndPointPlayerReady                = BaseUrl + "/match/%s/player/%s/ready"
	//Match
	EndPointMatch                  = BaseUrl + "/match"
	EndPointMatchState             = BaseUrl + "/match/%s/state"
	EndPointJoinMatch              = BaseUrl + "/match/%v/join/%v"
	EndPointOpenMatch              = BaseUrl + "/open"
	EndPointMatchCutDeck           = BaseUrl + "/match/%s/cut/%s"
	EndPointDeckByPlayerAndMatchId = BaseUrl + "/match/%s/player/%s/deck"
	//Account
	EndPointLogin = BaseUrl + "/account/login/%s"
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
