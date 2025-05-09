// Package services are the CLIs connection to the backend. If you look in server/route you'll see a coorsponding package for each service.
package services

import (
	"bytes"
	"context"
	"io"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	BaseURL                            = "http://localhost:1323/v1"
	EndPointKitty                      = BaseURL + "/match/%s/player/%s/to/%s/kitty"
	EndPointPlay                       = BaseURL + "/match/%s/player/%s/to/%s/play"
	EndPointPlayerByForMatchAndAccount = BaseURL + "/match/%s/account/%s"
	EndPointMatch                      = BaseURL + "/match/%s"
	EndPointOpenMatch                  = BaseURL + "/open"
	EndPointMatchState                 = BaseURL + "/match/%s/state"
	EndPointMatchCutDeck               = BaseURL + "/match/%s/cut/%s"
	EndPointDeckByPlayerAndMatchID     = BaseURL + "/match/%s/player/%s/deck"
	EndPointJoinMatch                  = BaseURL + "/match/%s/join/%s"
	EndPointLogin                      = BaseURL + "/account/login/%s"
)

func url(url string, method string, json string) tea.Msg {
	var buf *bytes.Buffer
	buf = bytes.NewBufferString(json)

	if json == "" {
		buf = bytes.NewBuffer(nil)
	}

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, method, url, buf)
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
