package store

import (
	"encoding/json"
	"fmt"

	"github.com/bardic/gocrib/queries/queries"
	"github.com/bardic/gocrib/vo"
	"github.com/labstack/echo/v4"
)

type MatchStore struct {
	Store
}

func (p *MatchStore) GetMatch(ctx echo.Context, matchID int) (*vo.Match, error) {
	m, err := p.q().GetMatchById(ctx.Request().Context(), matchID)

	defer p.Close()

	if err != nil {
		return nil, err
	}

	var match *vo.Match
	err = json.Unmarshal(m, &match)
	if err != nil {
		return nil, err
	}

	return match, nil
}

func (p *MatchStore) UpdateMatchState(ctx echo.Context, params queries.UpdateMatchStateParams) (*vo.Match, error) {
	match, err := p.q().UpdateMatchState(ctx.Request().Context(), params)
	if err != nil {
		return nil, err
	}

	defer p.Close()

	return &vo.Match{
		ID:            match.ID,
		Creationdate:  match.Creationdate,
		Privatematch:  match.Privatematch,
		Elorangemin:   match.Elorangemax,
		Elorangemax:   match.Elorangemax,
		Cutgamecardid: match.Cutgamecardid,
	}, nil
}

func (p *MatchStore) UpateMatchCurrentPlayerTurn(
	ctx echo.Context,
	params queries.UpateMatchCurrentPlayerTurnParams,
) error {
	err := p.q().UpateMatchCurrentPlayerTurn(ctx.Request().Context(), params)
	defer p.Close()
	if err != nil {
		return err
	}

	return nil
}

func (p *MatchStore) CreateMatch(ctx echo.Context, params queries.CreateMatchParams) (*vo.Match, error) {
	m, err := p.q().CreateMatch(ctx.Request().Context(), params)

	defer p.Close()

	if err != nil {
		return nil, err
	}

	return &vo.Match{
		ID: m.ID,
	}, nil
}

func (p *MatchStore) PlayerJoinMatch(ctx echo.Context, params queries.PlayerJoinMatchParams) (*vo.Match, error) {
	err := p.q().PlayerJoinMatch(ctx.Request().Context(), params)

	defer p.Close()

	if err != nil {
		return nil, err
	}

	return &vo.Match{}, nil
}

func (p *MatchStore) UpdateMatchCut(ctx echo.Context, params queries.UpdateMatchCutParams) error {
	err := p.q().UpdateMatchCut(ctx.Request().Context(), params)
	defer p.Close()

	if err != nil {
		return err
	}

	return nil
}

func (p *MatchStore) UpdatePlayerTurnOrder(ctx echo.Context, params queries.UpdatePlayerTurnOrderParams) error {
	err := p.q().UpdatePlayerTurnOrder(ctx.Request().Context(), params)
	defer p.Close()
	if err != nil {
		return err
	}

	return nil
}

func (p *MatchStore) UpdateDealerForMatch(ctx echo.Context, params queries.UpdateDealerForMatchParams) error {
	err := p.q().UpdateDealerForMatch(ctx.Request().Context(), params)
	defer p.Close()
	if err != nil {
		return err
	}
	return nil
}

func (p *MatchStore) GetOpenMatches(ctx echo.Context) ([]*vo.Match, error) {
	matchesData, err := p.q().GetOpenMatches(ctx.Request().Context(), queries.GamestateNew)
	if err != nil {
		return nil, err
	}

	matches := make([]*vo.Match, len(matchesData))
	for i, m := range matchesData {
		matches[i] = &vo.Match{
			ID: m.ID,
		}
	}

	return matches, nil
}

func (p *MatchStore) GetMatchState(ctx echo.Context, matchID int) (string, error) {
	matchState, err := p.q().GetMatchStateById(ctx.Request().Context(), matchID)
	if err != nil {
		return "", fmt.Errorf("[MATCH][GetMatchState] Can't find state: %w", err)
	}
	return string(matchState), nil
}

func (p *MatchStore) GetDeck(ctx echo.Context, matchID int) ([]*vo.Deck, error) {
	deck, err := p.q().GetDeckForMatchId(ctx.Request().Context(), matchID)

	defer p.Close()

	if err != nil {
		return nil, err
	}

	gameDecks := make([]*vo.Deck, len(deck))

	for i, d := range deck {
		gameDecks[i] = &vo.Deck{
			ID:             d.ID,
			Cutmatchcardid: d.Cutmatchcardid,
			Matchid:        d.Matchid,
			Cards:          []*vo.Card{}, // TODO fix this,
		}
	}

	return gameDecks, nil
}
