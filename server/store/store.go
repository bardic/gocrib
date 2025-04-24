package store

import (
	"github.com/bardic/gocrib/queries/queries"
	"github.com/labstack/echo/v4"
)

type Store struct {
	q *queries.Queries
	c echo.Context
}
