package match

import (
	"errors"
	"fmt"
	"net/http"

	logger "github.com/bardic/gocrib/cli/utils/log"
	"github.com/bardic/gocrib/queries/queries"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// GetOpenMatches route
//
//	@Summary	Get match by id
//	@Description
//	@Tags		match
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	[]vo.GameMatch
//	@Failure	404	{object}	error
//	@Failure	500	{object}	error
//	@Router		/open [get]
func (h *Handler) GetOpenMatches(c echo.Context) error {
	l := logger.Get()
	defer l.Sync()
	l.Log(zapcore.DebugLevel, "Get Open Matches")
	v, err := h.MatchStore.GetOpenMatches(c)
	if err != nil {
		e := errors.New(fmt.Sprintf("failed to get open matches: %v", err))
		l.Log(zapcore.ErrorLevel, e.Error())
		return c.JSON(http.StatusInternalServerError, e)
	}

	if v == nil {
		v = []queries.Match{}
	}

	l.Log(zapcore.DebugLevel, "Open Matches", zap.Field{
		Key:    "Match",
		Type:   zapcore.StringType,
		String: fmt.Sprintf("%v", v),
	})

	return c.JSON(http.StatusOK, v)
}
