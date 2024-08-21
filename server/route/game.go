package route

import (
	"net/http"

	"github.com/bardic/cribbage/server/model"
	"github.com/labstack/echo/v4"
)

// Create godoc
// @Summary      Play a card
// @Description
// @Tags         game
// @Accept       json
// @Produce      json
// @Param details body model.GameAction true "Action to perform"
// @Success      200  {object}  model.ScoreResults
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /game/playCard/ [post]
func PlayCard(c echo.Context) error {
	details := new(model.GameAction)
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	switch details.Type {
	case model.Cut:
		cutDeck(details.MatchId, details.Card)
	case model.Discard:
		discardCard(details.MatchId, details.Card)
	case model.Peg:
		pegCard(details.MatchId, details.Card)
	}

	return c.JSON(http.StatusOK, model.ScoreResults{Results: []model.Scores{
		{
			Cards: []int{0},
			Point: 0,
		},
	}})
}

func pegCard(matchId int, card model.GameplayCard) (model.ScoreResults, error) {
	UpdateCardsInPlay(matchId, card)

	res := model.ScoreResults{}

	r, err := scanForThirtyOne()
	if err != nil {
		return model.ScoreResults{}, err
	}
	res.Results = append(res.Results, r...)

	r, err = scanDeckCut()
	if err != nil {
		return model.ScoreResults{}, err
	}

	res.Results = append(res.Results, r...)
	if err != nil {
		return model.ScoreResults{}, err
	}
	res.Results = append(res.Results, r...)

	r, err = scanForFifthteens()
	if err != nil {
		return model.ScoreResults{}, err
	}
	res.Results = append(res.Results, r...)

	r, err = scanForLastCard()
	if err != nil {
		return model.ScoreResults{}, err
	}
	res.Results = append(res.Results, r...)

	r, err = scanForMatchingKinds()
	if err != nil {
		return model.ScoreResults{}, err
	}
	res.Results = append(res.Results, r...)

	r, err = scanForRuns()
	if err != nil {
		return model.ScoreResults{}, err
	}
	res.Results = append(res.Results, r...)

	return model.ScoreResults{Results: []model.Scores{
		{
			Cards: []int{0},
			Point: 0,
		},
	}}, nil
}

func discardCard(matchId int, card model.GameplayCard) (model.ScoreResults, error) {
	return model.ScoreResults{Results: []model.Scores{
		{
			Cards: []int{0},
			Point: 0,
		},
	}}, nil
}

func cutDeck(matchId int, card model.GameplayCard) (model.ScoreResults, error) {

	UpdateCut(matchId, card)

	return model.ScoreResults{Results: []model.Scores{
		{
			Cards: []int{0},
			Point: 0,
		},
	}}, nil
}

func scanForRuns() ([]model.Scores, error) {
	return []model.Scores{}, nil
}

func scanForMatchingKinds() ([]model.Scores, error) {
	return []model.Scores{}, nil
}

func scanForFifthteens() ([]model.Scores, error) {
	return []model.Scores{}, nil
}

func scanDeckCut() ([]model.Scores, error) {
	return []model.Scores{}, nil
}

func scanForThirtyOne() ([]model.Scores, error) {
	return []model.Scores{}, nil
}

func scanForLastCard() ([]model.Scores, error) {
	return []model.Scores{}, nil
}
