package match

import (
	"net/http"
	"strconv"

	"github.com/bardic/gocrib/queries/queries"
	"go.uber.org/zap/zapcore"

	"github.com/bardic/gocrib/vo"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"

	logger "github.com/bardic/gocrib/cli/utils/log"
)

var zero = 0

// NewMatch route
// @Summary	Create new match
// @Description
// @Tags		match
// @Accept		json
// @Produce	json
// @Param		accountId	path		int	true	"account id"'
// @Success	200		{object}	int
// @Failure	400		{object}	error
// @Failure	404		{object}	error
// @Failure	500		{object}	error
// @Router		/match/{accountId} [post]
func (h *Handler) NewMatch(c echo.Context) error {
	l := logger.Get()
	defer l.Sync()

	accountID, err := strconv.Atoi(c.Param("accountId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	deck, err := h.DeckStore.CreateDeck(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	match, err := h.MatchStore.CreateMatch(c, queries.CreateMatchParams{
		Privatematch:       false,
		Elorangemin:        &zero,
		Elorangemax:        &zero,
		Deckid:             deck.ID,
		Cutgamecardid:      &zero,
		Turnpasstimestamps: []pgtype.Timestamptz{},
		Gamestate:          queries.GamestateNew,
		Art:                "default.png",
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	cards, err := h.CardStore.GetCards(c)
	if err != nil {
		return err
	}

	l.Log(zapcore.DebugLevel, "Cards")

	deckCards := make([]queries.AddCardToDeckParams, len(cards))
	matchCards := make([]queries.CreateMatchCardParams, len(cards))

	for i, card := range cards {
		matchCards[i] = queries.CreateMatchCardParams{
			Cardid: card.ID,
			State:  queries.CardstateDeck,
		}

		deckCards[i] = queries.AddCardToDeckParams{
			Deckid:      deck.ID,
			Matchcardid: card.ID, // TODO: This is wrong. should get the id from the matchcards table, not the card id
		}
	}

	err = h.CardStore.CreateMatchCard(c, matchCards)
	if err != nil {
		return err
	}

	err = h.DeckStore.AddCardToDeck(c, deckCards)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	score := 0

	player, err := h.PlayerStore.CreatePlayer(c, queries.CreatePlayerParams{
		Accountid: &accountID,
		Score:     &score,
		Isready:   false,
		Art:       "default.png",
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	err = h.MatchStore.PlayerJoinMatch(c, queries.PlayerJoinMatchParams{
		Matchid:  match.ID,
		Playerid: player.ID,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, vo.MatchDetailsResponse{
		MatchID:   match.ID,
		PlayerID:  player.ID,
		GameState: match.Gamestate,
	})
}
