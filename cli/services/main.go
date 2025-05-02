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
	// Player
	EndPointKitty                      = BaseUrl + "/match/%s/player/%s/to/%s/kitty"
	EndPointPlay                       = BaseUrl + "/match/%s/player/%s/to/%s/play"
	EndPointPlayerByForMatchAndAccount = BaseUrl + "/match/%s/account/%s"
	// Match
	EndPointMatch                  = BaseUrl + "/match/%s"
	EndPointOpenMatch              = BaseUrl + "/open"
	EndPointMatchState             = BaseUrl + "/match/%s/state"
	EndPointMatchCutDeck           = BaseUrl + "/match/%s/cut/%s"
	EndPointDeckByPlayerAndMatchId = BaseUrl + "/match/%s/player/%s/deck"
	EndPointJoinMatch              = BaseUrl + "/match/%s/join/%s"
	// Account
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
