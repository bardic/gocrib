package route

import (
	"net/http"
	"strconv"

	"github.com/bardic/gocrib/server/store"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

// NewMatch route
//
//	@Summary			Create new match with accountId
//	@Description
//	@Tags					match
//	@Accept				json
//	@Produce			json
//	@Param				accountId		path			int		true	"account id"'
//	@Success			200					{object}	int
//	@Failure			400					{object}	error
//	@Failure			500					{object}	error
//	@Router				/match/{accountId} 	[post]
func (h *Handler) NewMatch(c echo.Context) error {
	accountID, err := strconv.Atoi(c.Param("accountId"))
	if err != nil {
		return h.BadParams(c, "error parsing accountID for NewMatch", err)
	}

	player, err := h.PlayerStore.CreatePlayer(c, accountID)
	if err != nil {
		return h.InternalError(c, "failed to create player", err)
	}

	zero := 0

	match, err := h.MatchStore.CreateMatch(c, store.NewMatchParam{
		Privatematch:       false,
		Elorangemin:        zero,
		Elorangemax:        zero,
		Cutgamecardid:      zero,
		Turnpasstimestamps: []pgtype.Timestamptz{},
		DealerID:           player.ID,
		CurrentPlayerID:    player.ID,
		Gamestate:          "'New'",
		Art:                "default.png",
	})
	if err != nil {
		return h.InternalError(c, "failed to create match in NewMatch", err)
	}

	deck, err := h.DeckStore.CreateDeck(c, match.ID)
	if err != nil {
		return h.InternalError(c, "", err)
	}

	cards, err := h.CardStore.GetCards(c)
	if err != nil {
		return h.InternalError(c, "failed to retreive base cards from DB", err)
	}
	cardIDs := make([]int, len(cards))

	for i, card := range cards {
		cardIDs[i] = card.ID
	}

	err = h.CardStore.CreateMatchCard(c, deck.ID, cardIDs)
	if err != nil {
		return h.InternalError(c, "failed to create match cards", err)
	}

	match, err = h.MatchStore.PlayerJoinMatch(c, match.ID, player.ID)
	if err != nil {
		return h.InternalError(c, "failed to have player join match", err)
	}

	return c.JSON(http.StatusOK, match)
}
