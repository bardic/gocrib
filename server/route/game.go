package route

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

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

	//Get all match details

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
			Cards: []model.GameplayCard{},
			Point: 0,
		},
	}})
}

func pegCard(matchId int, card model.GameplayCard) (model.ScoreResults, error) {
	match, err := UpdateCardsInPlay(matchId, card)
	if err != nil {
		return model.ScoreResults{}, err
	}
	res := model.ScoreResults{}

	r, err := scanForThirtyOne(match.CardsInPlay)
	if err != nil {
		return model.ScoreResults{}, err
	}
	res.Results = append(res.Results, r...)

	r, err = scanRightJackCut(match.CardsInPlay, match)
	if err != nil {
		return model.ScoreResults{}, err
	}

	r, err = scanJackOnCut(match)
	if err != nil {
		return model.ScoreResults{}, err
	}

	res.Results = append(res.Results, r...)
	if err != nil {
		return model.ScoreResults{}, err
	}
	res.Results = append(res.Results, r...)

	r, err = scanForFifthteens(match.CardsInPlay)
	if err != nil {
		return model.ScoreResults{}, err
	}
	res.Results = append(res.Results, r...)

	r, err = scanForLastCard(match)
	if err != nil {
		return model.ScoreResults{}, err
	}
	res.Results = append(res.Results, r...)

	// r, err = scanForMatchingKinds(match.CardsInPlay, []model.Scores{})
	// if err != nil {
	// 	return model.ScoreResults{}, err
	// }
	// res.Results = append(res.Results, r...)

	// r, err = scanForRuns(match.CardsInPlay)
	// if err != nil {
	// 	return model.ScoreResults{}, err
	// }
	// res.Results = append(res.Results, r...)

	return model.ScoreResults{Results: []model.Scores{
		{
			Cards: []model.GameplayCard{},
			Point: 0,
		},
	}}, nil
}

func discardCard(matchId int, card model.GameplayCard) (model.ScoreResults, error) {
	return model.ScoreResults{Results: []model.Scores{
		{
			Cards: []model.GameplayCard{},
			Point: 0,
		},
	}}, nil
}

func cutDeck(matchId int, card model.GameplayCard) (model.ScoreResults, error) {

	UpdateCut(matchId, card)

	return model.ScoreResults{Results: []model.Scores{
		{
			Cards: []model.GameplayCard{},
			Point: 0,
		},
	}}, nil
}

func scanForRuns(gameplayCardsIdsInPlay []int) ([]model.Scores, error) {
	cardsInPlay, err := getGameplayCardsForIds(gameplayCardsIdsInPlay)
	if err != nil {
		return []model.Scores{}, err
	}

	sort.Slice(cardsInPlay, func(i, j int) bool {
		return cardsInPlay[i].Value < cardsInPlay[j].Value
	})

	var pointsFound []model.Scores
	if cardsInPlay[0].Value+1 == cardsInPlay[1].Value &&
		cardsInPlay[1].Value+1 == cardsInPlay[2].Value {

		pointsFound = []model.Scores{{
			Cards: []model.GameplayCard{cardsInPlay[0], cardsInPlay[1], cardsInPlay[2]},
			Point: 3,
		}}

		if cardsInPlay[2].Value+1 == cardsInPlay[3].Value {
			pointsFound = []model.Scores{{
				Cards: []model.GameplayCard{cardsInPlay[0], cardsInPlay[1], cardsInPlay[2], cardsInPlay[3]},
				Point: 4,
			}}
		}
	}

	if cardsInPlay[0].Value+1 == cardsInPlay[1].Value &&
		cardsInPlay[1].Value+1 == cardsInPlay[3].Value {

		pointsFound = append(pointsFound, model.Scores{
			Cards: []model.GameplayCard{cardsInPlay[0], cardsInPlay[1], cardsInPlay[3]},
			Point: 3,
		})
	}

	return pointsFound, nil
}

func scanForMatchingKinds(gameplayCardsIdsInPlay []int) ([]model.Scores, error) {
	cardsInPlay, err := getGameplayCardsForIds(gameplayCardsIdsInPlay)
	if err != nil {
		return []model.Scores{}, err
	}

	sort.Slice(cardsInPlay, func(i, j int) bool {
		return cardsInPlay[i].Value < cardsInPlay[j].Value
	})

	var pointsFound []model.Scores

	if cardsInPlay[0].Value == cardsInPlay[1].Value {
		pointsFound = append(pointsFound, model.Scores{
			Cards: []model.GameplayCard{cardsInPlay[0], cardsInPlay[1]},
			Point: 2,
		})
	}

	if cardsInPlay[0].Value == cardsInPlay[2].Value {
		pointsFound = append(pointsFound, model.Scores{
			Cards: []model.GameplayCard{cardsInPlay[0], cardsInPlay[1]},
			Point: 2,
		})
	}

	if cardsInPlay[0].Value == cardsInPlay[3].Value {
		pointsFound = append(pointsFound, model.Scores{
			Cards: []model.GameplayCard{cardsInPlay[0], cardsInPlay[1]},
			Point: 2,
		})
	}

	if cardsInPlay[1].Value == cardsInPlay[2].Value {
		pointsFound = append(pointsFound, model.Scores{
			Cards: []model.GameplayCard{cardsInPlay[0], cardsInPlay[1]},
			Point: 2,
		})
	}

	if cardsInPlay[1].Value == cardsInPlay[3].Value {
		pointsFound = append(pointsFound, model.Scores{
			Cards: []model.GameplayCard{cardsInPlay[0], cardsInPlay[1]},
			Point: 2,
		})
	}

	if cardsInPlay[2].Value == cardsInPlay[3].Value {
		pointsFound = append(pointsFound, model.Scores{
			Cards: []model.GameplayCard{cardsInPlay[0], cardsInPlay[1]},
			Point: 2,
		})
	}

	return pointsFound, nil
}

func scanForFifthteens(gameplayCardsIdsInPlay []int) ([]model.Scores, error) {
	cardsInPlay, err := getGameplayCardsForIds(gameplayCardsIdsInPlay)
	if err != nil {
		return []model.Scores{}, err
	}

	var pointsFound []model.Scores

	for i := 0; i < len(cardsInPlay); i++ {
		for j := i; j < len(cardsInPlay); j++ {
			for k := j; k < len(cardsInPlay); k++ {
				if cardsInPlay[i].Value+cardsInPlay[j].Value+cardsInPlay[k].Value == 15 && (i != j && i != k && j != k) {
					pointsFound = append(pointsFound, model.Scores{
						Cards: []model.GameplayCard{cardsInPlay[0], cardsInPlay[1]},
						Point: 2,
					})
				}
			}
		}
	}

	for i := 0; i < len(cardsInPlay); i++ {
		for j := i; j < len(cardsInPlay); j++ {
			if cardsInPlay[i].Value+cardsInPlay[j].Value == 15 && i != j {
				pointsFound = append(pointsFound, model.Scores{
					Cards: []model.GameplayCard{cardsInPlay[0], cardsInPlay[1]},
					Point: 2,
				})
			}
		}
	}

	return pointsFound, echo.ErrHTTPVersionNotSupported.Internal
}

func scanRightJackCut(gameplayCardsIdsInPlay []int, match model.Match) ([]model.Scores, error) {
	cardsInPlay, err := getGameplayCardsForIds(gameplayCardsIdsInPlay)
	if err != nil {
		return []model.Scores{}, err
	}

	cut, err := getGameplayCardsForIds([]int{match.CutGameCardId})
	if err != nil {
		return []model.Scores{}, err
	}

	for i := 0; i < len(cardsInPlay); i++ {
		if cardsInPlay[i].Value == 11 && cardsInPlay[i].Suit == cut[0].Suit {
			return []model.Scores{{
				Cards: []model.GameplayCard{cardsInPlay[0], cardsInPlay[1]},
				Point: 1,
			}}, nil
		}
	}

	return []model.Scores{}, nil
}

func getGameplayCardsForIds(ids []int) ([]model.GameplayCard, error) {
	string_ids := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ids)), ","), "[]")
	cards, err := QueryForCards(string_ids)
	if err != nil {
		return []model.GameplayCard{}, err
	}
	return cards, nil
}

func scanForThirtyOne(gameplayCardsIdsInPlay []int) ([]model.Scores, error) {
	cardsInPlay, err := getGameplayCardsForIds(gameplayCardsIdsInPlay)
	if err != nil {
		return []model.Scores{}, err
	}

	plays := []model.GameplayCard{}
	pointsFound := []model.Scores{}

	if cardsInPlay[0].Value+cardsInPlay[1].Value+cardsInPlay[2].Value+cardsInPlay[3].Value == 31 {
		pointsFound = append(pointsFound, model.Scores{
			Cards: plays,
			Point: 2,
		})
	}

	if cardsInPlay[0].Value+cardsInPlay[1].Value+cardsInPlay[2].Value == 31 {
		pointsFound = append(pointsFound, model.Scores{
			Cards: plays,
			Point: 2,
		})
	}

	if cardsInPlay[0].Value+cardsInPlay[1].Value+cardsInPlay[3].Value == 31 {
		pointsFound = append(pointsFound, model.Scores{
			Cards: plays,
			Point: 2,
		})
	}

	if cardsInPlay[0].Value+cardsInPlay[2].Value+cardsInPlay[3].Value == 31 {
		pointsFound = append(pointsFound, model.Scores{
			Cards: plays,
			Point: 2,
		})
	}

	if cardsInPlay[1].Value+cardsInPlay[2].Value+cardsInPlay[3].Value == 31 {
		pointsFound = append(pointsFound, model.Scores{
			Cards: plays,
			Point: 2,
		})
	}

	return pointsFound, nil
}

func scanForLastCard(match model.Match) ([]model.Scores, error) {

	// playerOne, err := getGameplayCardsForIds([]int{match.PlayerId[0]})
	// playerTwo, err := getGameplayCardsForIds([]int{match.CutGameCardId})

	return []model.Scores{}, nil
}

func scanJackOnCut(match model.Match) ([]model.Scores, error) {
	cardsInPlay, err := getGameplayCardsForIds([]int{match.CutGameCardId})
	if err != nil {
		return []model.Scores{}, err
	}

	if cardsInPlay[0].Value == 11 {
		return []model.Scores{{
			Cards: []model.GameplayCard{},
			Point: 2,
		}}, nil
	}

	return []model.Scores{}, nil
}
