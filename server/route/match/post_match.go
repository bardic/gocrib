package match

import (
	"net/http"
	"strconv"

	"github.com/bardic/gocrib/queries/queries"

	"github.com/bardic/gocrib/vo"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

var (
	Zero = 0
)

// Create godoc
//
//	@Summary	Create new match
//	@Description
//	@Tags		match
//	@Accept		json
//	@Produce	json
//	@Param		accountId	path		int	true	"account id"'
//	@Success	200		{object}	int
//	@Failure	400		{object}	error
//	@Failure	404		{object}	error
//	@Failure	500		{object}	error
//	@Router		/match/{accountId} [post]
func (h *MatchHandler) NewMatch(c echo.Context) error {
	accountId, err := strconv.Atoi(c.Param("accountId"))

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	deck, err := h.DeckStore.CreateDeck(c)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	match, err := h.MatchStore.CreateMatch(c, queries.CreateMatchParams{
		Privatematch:       false,
		Elorangemin:        &Zero,
		Elorangemax:        &Zero,
		Deckid:             deck.ID,
		Cutgamecardid:      &Zero,
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

	for _, card := range cards {
		matchCard, err := h.CardStore.CreateMatchCard(c, queries.CreateMatchCardParams{
			Cardid: card.ID,
			State:  queries.CardstateDeck,
		})

		if err != nil {
			return err
		}

		err = h.DeckStore.AddCardToDeck(c, queries.AddCardToDeckParams{
			Deckid:      deck.ID,
			Matchcardid: matchCard.ID,
		})

		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
	}

	score := 0

	player, err := h.PlayerStore.CreatePlayer(c, queries.CreatePlayerParams{
		Accountid: &accountId,
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
		MatchId:   match.ID,
		PlayerId:  player.ID,
		GameState: match.Gamestate,
	})
}
