package game

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/bardic/gocrib/model"
	"github.com/bardic/gocrib/server/utils"
)

func cardsInPlay(players []model.Player) []int {
	cardIds := []int{}
	for _, p := range players {
		cardIds = append(cardIds, p.Hand...)
	}

	return cardIds
}

func countPegs(m model.GameMatch) (model.ScoreResults, error) {
	res := model.ScoreResults{}

	r, err := scanForThirtyOne(cardsInPlay(m.Players))
	if err != nil {
		return model.ScoreResults{}, err
	}
	res.Results = append(res.Results, r...)

	if len(m.TurnPassTimestamps) == 0 {
		r, err = scanJackOnCut(m)
		if err != nil {
			return model.ScoreResults{}, err
		}
		res.Results = append(res.Results, r...)
	}

	r, err = scanForFifthteens(cardsInPlay(m.Players))
	if err != nil {
		return model.ScoreResults{}, err
	}
	res.Results = append(res.Results, r...)

	r, err = scanForLastCard(m)
	if err != nil {
		return model.ScoreResults{}, err
	}
	res.Results = append(res.Results, r...)

	r, err = scanForMatchingKinds(cardsInPlay(m.Players))
	if err != nil {
		return model.ScoreResults{}, err
	}
	res.Results = append(res.Results, r...)

	r, err = scanForRuns(cardsInPlay(m.Players))
	if err != nil {
		return model.ScoreResults{}, err
	}
	res.Results = append(res.Results, r...)

	return model.ScoreResults{Results: res.Results}, nil
}

func countHand(m model.GameMatch) (model.ScoreResults, error) {
	res := model.ScoreResults{}

	r, err := scanForThirtyOne(cardsInPlay(m.Players))
	if err != nil {
		return model.ScoreResults{}, err
	}
	res.Results = append(res.Results, r...)

	r, err = scanRightJackCut(cardsInPlay(m.Players), m)
	if err != nil {
		return model.ScoreResults{}, err
	}
	res.Results = append(res.Results, r...)

	r, err = scanForFifthteens(cardsInPlay(m.Players))
	if err != nil {
		return model.ScoreResults{}, err
	}
	res.Results = append(res.Results, r...)

	r, err = scanForLastCard(m)
	if err != nil {
		return model.ScoreResults{}, err
	}
	res.Results = append(res.Results, r...)

	r, err = scanForMatchingKinds(cardsInPlay(m.Players))
	if err != nil {
		return model.ScoreResults{}, err
	}
	res.Results = append(res.Results, r...)

	r, err = scanForRuns(cardsInPlay(m.Players))
	if err != nil {
		return model.ScoreResults{}, err
	}
	res.Results = append(res.Results, r...)

	r, err = scanForFlush(cardsInPlay(m.Players))
	if err != nil {
		return model.ScoreResults{}, err
	}
	res.Results = append(res.Results, r...)

	return model.ScoreResults{Results: res.Results}, nil
}

func scanForFlush(cardIdsInHand []int) ([]model.Scores, error) {
	gameplayCardsInHand, err := getGameplayCardsForIds(cardIdsInHand)
	if err != nil {
		return []model.Scores{}, err
	}

	var flushSuit *model.Suit
	matchesFlush := true
	for _, card := range gameplayCardsInHand {
		if flushSuit == nil {
			flushSuit = &card.Suit
		} else {
			matchesFlush = flushSuit == &card.Suit
		}
	}

	if matchesFlush {
		return []model.Scores{
			{
				Point: 4,
			},
		}, nil
	}

	return []model.Scores{}, nil
}

func discardCard(match model.GameMatch) (model.ScoreResults, error) {
	return model.ScoreResults{Results: []model.Scores{
		{
			Cards: []model.GameplayCard{},
			Point: 0,
		},
	}}, nil
}

func cutDeck(m model.GameMatch, cutCardId int) (model.ScoreResults, error) {

	utils.UpdateCut(m.Id, cutCardId)

	return model.ScoreResults{Results: []model.Scores{
		{
			Cards: []model.GameplayCard{},
			Point: 0,
		},
	}}, nil
}

func scanForRuns(cardIdsInPlay []int) ([]model.Scores, error) {
	gameplayCardsInPlay, err := getGameplayCardsForIds(cardIdsInPlay)
	if err != nil {
		return []model.Scores{}, err
	}

	sort.Slice(gameplayCardsInPlay, func(i, j int) bool {
		return gameplayCardsInPlay[i].Value < gameplayCardsInPlay[j].Value
	})

	var pointsFound []model.Scores

	if len(gameplayCardsInPlay) < 3 {
		return pointsFound, nil
	}

	if gameplayCardsInPlay[0].Value+1 == gameplayCardsInPlay[1].Value &&
		gameplayCardsInPlay[1].Value+1 == gameplayCardsInPlay[2].Value {

		pointsFound = []model.Scores{{
			Cards: []model.GameplayCard{gameplayCardsInPlay[0], gameplayCardsInPlay[1], gameplayCardsInPlay[2]},
			Point: 3,
		}}

		if len(gameplayCardsInPlay) < 4 {
			return pointsFound, nil
		}

		if gameplayCardsInPlay[2].Value+1 == gameplayCardsInPlay[3].Value {
			pointsFound = []model.Scores{{
				Cards: []model.GameplayCard{gameplayCardsInPlay[0], gameplayCardsInPlay[1], gameplayCardsInPlay[2], gameplayCardsInPlay[3]},
				Point: 4,
			}}
		}
	}

	if gameplayCardsInPlay[0].Value+1 == gameplayCardsInPlay[1].Value &&
		gameplayCardsInPlay[1].Value+1 == gameplayCardsInPlay[3].Value {

		pointsFound = append(pointsFound, model.Scores{
			Cards: []model.GameplayCard{gameplayCardsInPlay[0], gameplayCardsInPlay[1], gameplayCardsInPlay[3]},
			Point: 3,
		})
	}

	return pointsFound, nil
}

func scanForMatchingKinds(cardIdsInPlay []int) ([]model.Scores, error) {
	gameplayCardsInPlay, err := getGameplayCardsForIds(cardIdsInPlay)
	if err != nil {
		return []model.Scores{}, err
	}

	sort.Slice(gameplayCardsInPlay, func(i, j int) bool {
		return gameplayCardsInPlay[i].Value < gameplayCardsInPlay[j].Value
	})

	var pointsFound []model.Scores

	if len(gameplayCardsInPlay) < 2 {
		return pointsFound, nil
	}

	if gameplayCardsInPlay[0].Value == gameplayCardsInPlay[1].Value {
		pointsFound = append(pointsFound, model.Scores{
			Cards: []model.GameplayCard{gameplayCardsInPlay[0], gameplayCardsInPlay[1]},
			Point: 2,
		})
	}

	if len(gameplayCardsInPlay) < 3 {
		return pointsFound, nil
	}

	if gameplayCardsInPlay[0].Value == gameplayCardsInPlay[2].Value {
		pointsFound = append(pointsFound, model.Scores{
			Cards: []model.GameplayCard{gameplayCardsInPlay[0], gameplayCardsInPlay[2]},
			Point: 2,
		})
	}

	if gameplayCardsInPlay[1].Value == gameplayCardsInPlay[2].Value {
		pointsFound = append(pointsFound, model.Scores{
			Cards: []model.GameplayCard{gameplayCardsInPlay[1], gameplayCardsInPlay[2]},
			Point: 2,
		})
	}

	if len(gameplayCardsInPlay) < 4 {
		return pointsFound, nil
	}

	if gameplayCardsInPlay[0].Value == gameplayCardsInPlay[3].Value {
		pointsFound = append(pointsFound, model.Scores{
			Cards: []model.GameplayCard{gameplayCardsInPlay[0], gameplayCardsInPlay[3]},
			Point: 2,
		})
	}

	if gameplayCardsInPlay[1].Value == gameplayCardsInPlay[3].Value {
		pointsFound = append(pointsFound, model.Scores{
			Cards: []model.GameplayCard{gameplayCardsInPlay[1], gameplayCardsInPlay[3]},
			Point: 2,
		})
	}

	if gameplayCardsInPlay[2].Value == gameplayCardsInPlay[3].Value {
		pointsFound = append(pointsFound, model.Scores{
			Cards: []model.GameplayCard{gameplayCardsInPlay[2], gameplayCardsInPlay[3]},
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

	pointsFound := []model.Scores{}

	//find if any combination of cardsInPlay equals 15
	for i := 0; i < len(cardsInPlay); i++ {
		for j := i; j < len(cardsInPlay); j++ {
			if cardsInPlay[i].Value+cardsInPlay[j].Value == 15 {
				pointsFound = append(pointsFound, model.Scores{
					Cards: []model.GameplayCard{cardsInPlay[i], cardsInPlay[j]},
					Point: 2,
				})
			}
		}
	}

	if len(cardsInPlay) < 3 {
		return pointsFound, nil
	}

	for i := 0; i < len(cardsInPlay); i++ {
		for j := i; j < len(cardsInPlay); j++ {
			for k := j; k < len(cardsInPlay); k++ {
				if cardsInPlay[i].Value+cardsInPlay[j].Value+cardsInPlay[k].Value == 15 && i != j && j != k {
					pointsFound = append(pointsFound, model.Scores{
						Cards: []model.GameplayCard{cardsInPlay[i], cardsInPlay[j], cardsInPlay[k]},
						Point: 2,
					})
				}
			}
		}
	}

	return pointsFound, nil
}

func scanRightJackCut(gameplayCardsIdsInPlay []int, match model.GameMatch) ([]model.Scores, error) {
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
	if len(ids) == 0 {
		return []model.GameplayCard{}, nil
	}

	string_ids := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ids)), ","), "[]")
	cards, err := utils.QueryForCards(string_ids)
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

	if len(cardsInPlay) < 3 {
		return pointsFound, nil
	}

	if cardsInPlay[0].Value+cardsInPlay[1].Value+cardsInPlay[2].Value == 31 {
		pointsFound = append(pointsFound, model.Scores{
			Cards: plays,
			Point: 2,
		})
	}

	if len(cardsInPlay) < 4 {
		return pointsFound, nil
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

	if cardsInPlay[0].Value+cardsInPlay[1].Value+cardsInPlay[2].Value+cardsInPlay[3].Value == 31 {
		pointsFound = append(pointsFound, model.Scores{
			Cards: plays,
			Point: 2,
		})
	}

	return pointsFound, nil
}

func scanForLastCard(m model.GameMatch) ([]model.Scores, error) {
	if len(m.Players) < 2 {
		return []model.Scores{}, errors.New("too few players")
	}

	playerOneCards, err := getGameplayCardsForIds(m.Players[0].Hand)

	if err != nil {
		return []model.Scores{}, err
	}

	playerTwoCards, err := getGameplayCardsForIds(m.Players[1].Hand)

	if err != nil {
		return []model.Scores{}, err
	}

	cardsInPlay, err := getGameplayCardsForIds(cardsInPlay(m.Players))

	if err != nil {
		return []model.Scores{}, err
	}

	total := 0
	for _, card := range cardsInPlay {
		total = total + int(card.Value)
	}

	playerOneCanPlay := false
	for _, card := range playerOneCards {
		if int(card.Value)+total <= 31 {
			playerOneCanPlay = true
		}
	}

	playerTwoCanPlay := false
	for _, card := range playerTwoCards {
		if int(card.Value)+total <= 31 {
			playerTwoCanPlay = true
		}
	}

	if !playerOneCanPlay && !playerTwoCanPlay {
		return []model.Scores{
			{
				Point: 1,
			},
		}, nil
	}

	return []model.Scores{}, nil
}

func scanJackOnCut(match model.GameMatch) ([]model.Scores, error) {
	cardsInPlay, err := getGameplayCardsForIds([]int{match.CutGameCardId})
	if err != nil || len(cardsInPlay) == 0 {
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
