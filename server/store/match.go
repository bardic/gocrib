package store

import (
	"encoding/json"
	"fmt"

	"github.com/bardic/gocrib/queries/queries"
	"github.com/bardic/gocrib/vo"
	"github.com/jackc/pgx/v5/pgtype"
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

func (p *MatchStore) UpdateMatchState(ctx echo.Context, matchID int, matchState string) (*vo.Match, error) {
	match, err := p.q().UpdateMatchState(ctx.Request().Context(), queries.UpdateMatchStateParams{
		Gamestate: matchState,
		ID:        matchID,
	})
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

func (p *MatchStore) UpateMatchCurrentPlayerTurn(ctx echo.Context, matchID, currentPlayereturnID int) error {
	err := p.q().UpateMatchCurrentPlayerTurn(ctx.Request().Context(), queries.UpateMatchCurrentPlayerTurnParams{
		Currentplayerturn: currentPlayereturnID,
		ID:                matchID,
	})
	defer p.Close()
	if err != nil {
		return err
	}

	return nil
}

type NewMatchParam struct {
	Privatematch       bool
	Elorangemin        int
	Elorangemax        int
	Cutgamecardid      int
	Turnpasstimestamps []pgtype.Timestamptz
	Gamestate          string
	DealerID           int
	CurrentPlayerID    int
	Art                string
}

func (p *MatchStore) CreateMatch(ctx echo.Context, params NewMatchParam) (*vo.Match, error) {
	m, err := p.q().CreateMatch(ctx.Request().Context(), queries.CreateMatchParams{
		Privatematch:       params.Privatematch,
		Elorangemin:        params.Elorangemin,
		Elorangemax:        params.Elorangemax,
		Cutgamecardid:      params.Cutgamecardid,
		Turnpasstimestamps: params.Turnpasstimestamps,
		Dealerid:           params.DealerID,
		Currentplayerturn:  params.CurrentPlayerID,
		Gamestate:          "'New'",
		Art:                params.Art,
	})

	defer p.Close()

	if err != nil {
		return nil, err
	}

	return &vo.Match{
		ID: m.ID,
	}, nil
}

func (p *MatchStore) PlayerJoinMatch(ctx echo.Context, matchID, playerID int) (*vo.Match, error) {
	err := p.q().PlayerJoinMatch(ctx.Request().Context(), queries.PlayerJoinMatchParams{
		Matchid:  matchID,
		Playerid: playerID,
	})

	defer p.Close()

	if err != nil {
		return nil, err
	}

	return &vo.Match{}, nil
}

func (p *MatchStore) UpdateMatchCut(ctx echo.Context, matchID, cutIndex int) error {
	err := p.q().UpdateMatchCut(ctx.Request().Context(), queries.UpdateMatchCutParams{
		ID:            matchID,
		Cutgamecardid: cutIndex,
	})
	defer p.Close()

	if err != nil {
		return err
	}

	return nil
}

func (p *MatchStore) UpdatePlayerTurnOrder(ctx echo.Context, matchID, playerID, turnOrder int) error {
	err := p.q().UpdatePlayerTurnOrder(ctx.Request().Context(), queries.UpdatePlayerTurnOrderParams{
		Turnorder: turnOrder,
		Matchid:   matchID,
		Playerid:  playerID,
	})
	defer p.Close()
	if err != nil {
		return err
	}

	return nil
}

func (p *MatchStore) UpdateDealerForMatch(ctx echo.Context, matchID, dealerID int) error {
	err := p.q().UpdateDealerForMatch(ctx.Request().Context(), queries.UpdateDealerForMatchParams{
		Dealerid: dealerID,
		ID:       matchID,
	})
	defer p.Close()
	if err != nil {
		return err
	}
	return nil
}

func (p *MatchStore) GetOpenMatches(ctx echo.Context) ([]*vo.Match, error) {
	matchesData, err := p.q().GetOpenMatches(ctx.Request().Context(), "'New'")
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

func (p *MatchStore) SetMatchState(ctx echo.Context, matchID int, state string) error {
	_, err := p.q().UpdateMatchState(
		ctx.Request().Context(),
		queries.UpdateMatchStateParams{
			Gamestate: state,
			ID:        matchID,
		},
	)
	if err != nil {
		return err
	}

	return nil
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
