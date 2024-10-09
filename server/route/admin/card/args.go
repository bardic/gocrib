package card

import (
	"github.com/bardic/gocrib/model"
	"github.com/jackc/pgx/v5"
)

func parseCard(details model.Card) pgx.NamedArgs {
	return pgx.NamedArgs{
		"id":    details.Id,
		"value": details.Value,
		"suit":  details.Suit,
		"art":   details.Art,
	}
}
