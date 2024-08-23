package route

import (
	"net/http"
	"slices"

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
	match, err := UpdateCardsInPlay(matchId, card)

	res := model.ScoreResults{}

	r, err := scanForThirtyOne(match.CardsInPlay)
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

	r, err = scanForMatchingKinds(match.CardsInPlay, []model.Scores{})
	if err != nil {
		return model.ScoreResults{}, err
	}
	res.Results = append(res.Results, r...)

	r, err = scanForRuns(match.CardsInPlay)
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

func scanForRuns(cardsInPlay []int) ([]model.Scores, error) {

	slices.Sort(cardsInPlay)
	var plays []int
	pointsFound := []model.Scores{}
	if cardsInPlay[0]+1 == cardsInPlay[1] &&
		cardsInPlay[1]+1 == cardsInPlay[2] {
		plays = []int{cardsInPlay[0], cardsInPlay[1], cardsInPlay[2]}
		pointsFound = append(pointsFound, model.Scores{
			Cards: plays,
			Point: 3,
		})

		if cardsInPlay[2]+1 == cardsInPlay[3] {
			plays = append(cardsInPlay, cardsInPlay[3])
			pointsFound = append(pointsFound, model.Scores{
				Cards: plays,
				Point: 4,
			})
		}
	}

	return pointsFound, nil
}

func scanForMatchingKinds(cardsInPlay []int, pointsFound []model.Scores) ([]model.Scores, error) {
	slices.Sort(cardsInPlay)

	if len(cardsInPlay) > 0 && cardsInPlay[0] == cardsInPlay[1] {
		plays := []int{cardsInPlay[0], cardsInPlay[1]}
		pointsFound = append(pointsFound, model.Scores{
			Cards: plays,
			Point: 2,
		})
		if len(cardsInPlay) > 1 && cardsInPlay[1] == cardsInPlay[2] {
			plays = append(plays, cardsInPlay[2])
			pointsFound = append(pointsFound, model.Scores{
				Cards: plays,
				Point: 6,
			})
			if len(cardsInPlay) > 2 && cardsInPlay[2] == cardsInPlay[3] {
				plays = append(plays, cardsInPlay[3])
				pointsFound = append(pointsFound, model.Scores{
					Cards: plays,
					Point: 12,
				})
			}
		}
	}

	if len(cardsInPlay) > 2 {
		scanForMatchingKinds(cardsInPlay[1:], pointsFound)
	}

	return pointsFound, nil
}

func scanForFifthteens() ([]model.Scores, error) {
	return []model.Scores{}, nil
}

func scanDeckCut() ([]model.Scores, error) {
	return []model.Scores{}, nil
}

func scanForThirtyOne(cardsInPlay []int) ([]model.Scores, error) {
	count := 0
	plays := []int{}
	pointsFound := []model.Scores{}

	for _, card := range cardsInPlay {
		count += card //todo fix this
		plays = append(plays, card)
		if count == 31 {
			pointsFound = append(pointsFound, model.Scores{
				Cards: plays,
				Point: 2,
			})

			if len(plays) < len(cardsInPlay) {
				scanForThirtyOne(cardsInPlay[1:])
			}
		}
	}

	return pointsFound, nil
}

func scanForLastCard() ([]model.Scores, error) {
	return []model.Scores{}, nil
}
