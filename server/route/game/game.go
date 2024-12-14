package game

import (
	"errors"
	"queries"
	"sort"

	"github.com/bardic/gocrib/vo"
)

func cardsInPlay(players []*queries.Player) []int32 {
	cardIds := []int32{}
	for _, p := range players {
		cardIds = append(cardIds, p.Hand...)
	}

	return cardIds
}

func countPegs(m vo.GameMatch) (vo.ScoreResults, error) {
	res := vo.ScoreResults{}

	r, err := scanForThirtyOne(cardsInPlay(m.Players))
	if err != nil {
		return vo.ScoreResults{}, err
	}
	res.Results = append(res.Results, r...)

	if len(m.Turnpasstimestamps) == 0 {
		r, err = scanJackOnCut(m)
		if err != nil {
			return vo.ScoreResults{}, err
		}
		res.Results = append(res.Results, r...)
	}

	r, err = scanForFifthteens(cardsInPlay(m.Players))
	if err != nil {
		return vo.ScoreResults{}, err
	}
	res.Results = append(res.Results, r...)

	r, err = scanForLastCard(m)
	if err != nil {
		return vo.ScoreResults{}, err
	}
	res.Results = append(res.Results, r...)

	r, err = scanForMatchingKinds(cardsInPlay(m.Players))
	if err != nil {
		return vo.ScoreResults{}, err
	}
	res.Results = append(res.Results, r...)

	r, err = scanForRuns(cardsInPlay(m.Players))
	if err != nil {
		return vo.ScoreResults{}, err
	}
	res.Results = append(res.Results, r...)

	return vo.ScoreResults{Results: res.Results}, nil
}

func countHand(m vo.GameMatch) (vo.ScoreResults, error) {
	res := vo.ScoreResults{}

	r, err := scanForThirtyOne(cardsInPlay(m.Players))
	if err != nil {
		return vo.ScoreResults{}, err
	}
	res.Results = append(res.Results, r...)

	r, err = scanRightJackCut(cardsInPlay(m.Players), m)
	if err != nil {
		return vo.ScoreResults{}, err
	}
	res.Results = append(res.Results, r...)

	r, err = scanForFifthteens(cardsInPlay(m.Players))
	if err != nil {
		return vo.ScoreResults{}, err
	}
	res.Results = append(res.Results, r...)

	r, err = scanForLastCard(m)
	if err != nil {
		return vo.ScoreResults{}, err
	}
	res.Results = append(res.Results, r...)

	r, err = scanForMatchingKinds(cardsInPlay(m.Players))
	if err != nil {
		return vo.ScoreResults{}, err
	}
	res.Results = append(res.Results, r...)

	r, err = scanForRuns(cardsInPlay(m.Players))
	if err != nil {
		return vo.ScoreResults{}, err
	}
	res.Results = append(res.Results, r...)

	r, err = scanForFlush(cardsInPlay(m.Players))
	if err != nil {
		return vo.ScoreResults{}, err
	}
	res.Results = append(res.Results, r...)

	return vo.ScoreResults{Results: res.Results}, nil
}

func scanForFlush(cardIdsInHand []int32) ([]vo.Scores, error) {
	gameplayCardsInHand, err := getGameplayCardsForIds(cardIdsInHand)
	if err != nil {
		return []vo.Scores{}, err
	}

	var flushSuit *queries.Cardsuit
	matchesFlush := true
	for _, card := range gameplayCardsInHand {
		if flushSuit == nil {
			flushSuit = &card.Suit
		} else {
			matchesFlush = flushSuit == &card.Suit
		}
	}

	if matchesFlush {
		return []vo.Scores{
			{
				Point: 4,
			},
		}, nil
	}

	return []vo.Scores{}, nil
}

func scanForRuns(cardIdsInPlay []int32) ([]vo.Scores, error) {
	gameplayCardsInPlay, err := getGameplayCardsForIds(cardIdsInPlay)
	if err != nil {
		return []vo.Scores{}, err
	}

	sort.Slice(gameplayCardsInPlay, func(i, j int) bool {
		return gameplayCardsInPlay[i].Value < gameplayCardsInPlay[j].Value
	})

	var pointsFound []vo.Scores

	if len(gameplayCardsInPlay) < 3 {
		return pointsFound, nil
	}

	gameCard1 := gameplayCardsInPlay[0]
	gameCard2 := gameplayCardsInPlay[1]
	gameCard3 := gameplayCardsInPlay[2]

	details1 := cardDetails(gameCard1.Value)
	details2 := cardDetails(gameCard2.Value)
	details3 := cardDetails(gameCard3.Value)

	if details1.Value+1 == details2.Value && details2.Value+1 == details3.Value {

		pointsFound = []vo.Scores{{
			Cards: []vo.GameCard{gameCard1, gameCard2, gameCard3},
			Point: 3,
		}}

		if len(gameplayCardsInPlay) < 4 {
			return pointsFound, nil
		}

		gameCard4 := gameplayCardsInPlay[3]
		details4 := cardDetails(gameCard4.Value)

		if details3.Value+1 == details4.Value {
			pointsFound = []vo.Scores{{
				Cards: []vo.GameCard{gameCard1, gameCard2, gameCard3, gameCard4},
				Point: 4,
			}}
		}
	}

	if details1.Value+1 == details2.Value &&
		details2.Value+1 == details3.Value {

		pointsFound = append(pointsFound, vo.Scores{
			Cards: []vo.GameCard{gameCard1, gameCard2, gameCard3},
			Point: 3,
		})
	}

	return pointsFound, nil
}

func cardDetails(cardValue queries.Cardvalue) vo.GameCardDetails {
	switch cardValue {
	case queries.CardvalueAce:

		return vo.GameCardDetails{
			Value: 1,
			Order: 1,
		}
	case queries.CardvalueTwo:
		return vo.GameCardDetails{
			Value: 2,
			Order: 2,
		}
	case queries.CardvalueThree:
		return vo.GameCardDetails{
			Value: 3,
			Order: 3,
		}
	case queries.CardvalueFour:
		return vo.GameCardDetails{
			Value: 4,
			Order: 4,
		}
	case queries.CardvalueFive:
		return vo.GameCardDetails{
			Value: 5,
			Order: 5,
		}
	case queries.CardvalueSix:
		return vo.GameCardDetails{
			Value: 6,
			Order: 6,
		}
	case queries.CardvalueSeven:
		return vo.GameCardDetails{
			Value: 7,
			Order: 7,
		}
	case queries.CardvalueEight:
		return vo.GameCardDetails{
			Value: 8,
			Order: 8,
		}
	case queries.CardvalueNine:
		return vo.GameCardDetails{
			Value: 9,
			Order: 9,
		}
	case queries.CardvalueTen:
		return vo.GameCardDetails{
			Value: 10,
			Order: 10,
		}
	case queries.CardvalueJack:
		return vo.GameCardDetails{
			Value: 10,
			Order: 11,
		}
	case queries.CardvalueQueen:
		return vo.GameCardDetails{
			Value: 10,
			Order: 12,
		}
	case queries.CardvalueKing:
		return vo.GameCardDetails{
			Value: 10,
			Order: 13,
		}
	case queries.CardvalueJoker:
		return vo.GameCardDetails{
			Value: 0,
			Order: 0,
		}
	}

	return vo.GameCardDetails{
		Value: 0,
		Order: 0,
	}

}

func scanForMatchingKinds(cardIdsInPlay []int32) ([]vo.Scores, error) {
	gameplayCardsInPlay, err := getGameplayCardsForIds(cardIdsInPlay)
	if err != nil {
		return []vo.Scores{}, err
	}

	sort.Slice(gameplayCardsInPlay, func(i, j int) bool {
		return gameplayCardsInPlay[i].Value < gameplayCardsInPlay[j].Value
	})

	var pointsFound []vo.Scores

	if len(gameplayCardsInPlay) < 2 {
		return pointsFound, nil
	}

	gameCard1 := gameplayCardsInPlay[0]
	gameCard2 := gameplayCardsInPlay[1]

	details1 := cardDetails(gameCard1.Value)
	details2 := cardDetails(gameCard2.Value)

	if details1.Value == details2.Value {
		pointsFound = append(pointsFound, vo.Scores{
			Cards: []vo.GameCard{gameCard1, gameCard2},
			Point: 2,
		})
	}

	if len(gameplayCardsInPlay) < 3 {
		return pointsFound, nil
	}

	gameCard3 := gameplayCardsInPlay[2]
	details3 := cardDetails(gameCard3.Value)

	if details1.Value == details3.Value {
		pointsFound = append(pointsFound, vo.Scores{
			Cards: []vo.GameCard{gameCard1, gameCard3},
			Point: 2,
		})
	}

	if details2.Value == details3.Value {
		pointsFound = append(pointsFound, vo.Scores{
			Cards: []vo.GameCard{gameCard2, gameCard3},
			Point: 2,
		})
	}

	if len(gameplayCardsInPlay) < 4 {
		return pointsFound, nil
	}

	gameCard4 := gameplayCardsInPlay[3]
	details4 := cardDetails(gameCard4.Value)

	if details1.Value == details4.Value {
		pointsFound = append(pointsFound, vo.Scores{
			Cards: []vo.GameCard{gameCard1, gameCard4},
			Point: 2,
		})
	}

	if details2.Value == details4.Value {
		pointsFound = append(pointsFound, vo.Scores{
			Cards: []vo.GameCard{gameCard2, gameCard4},
			Point: 2,
		})
	}

	if details3.Value == details4.Value {
		pointsFound = append(pointsFound, vo.Scores{
			Cards: []vo.GameCard{gameCard3, gameCard4},
			Point: 2,
		})
	}

	return pointsFound, nil
}

func scanForFifthteens(gameplayCardsIdsInPlay []int32) ([]vo.Scores, error) {
	cardsInPlay, err := getGameplayCardsForIds(gameplayCardsIdsInPlay)
	if err != nil {
		return []vo.Scores{}, err
	}

	pointsFound := []vo.Scores{}

	//find if any combination of cardsInPlay equals 15
	for i := 0; i < len(cardsInPlay); i++ {
		for j := i; j < len(cardsInPlay); j++ {
			details1 := cardDetails(cardsInPlay[i].Value)
			details2 := cardDetails(cardsInPlay[j].Value)
			// details3 := cardDetails(cardsInPlay[k].Value)

			if details1.Value+details2.Value == 15 {
				pointsFound = append(pointsFound, vo.Scores{
					Cards: []vo.GameCard{cardsInPlay[i], cardsInPlay[j]},
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
				details1 := cardDetails(cardsInPlay[i].Value)
				details2 := cardDetails(cardsInPlay[j].Value)
				details3 := cardDetails(cardsInPlay[k].Value)

				if details1.Value+details2.Value+details3.Value == 15 && i != j && j != k {
					pointsFound = append(pointsFound, vo.Scores{
						Cards: []vo.GameCard{cardsInPlay[i], cardsInPlay[j], cardsInPlay[k]},
						Point: 2,
					})
				}
			}
		}
	}

	return pointsFound, nil
}

func scanRightJackCut(gameplayCardsIdsInPlay []int32, match vo.GameMatch) ([]vo.Scores, error) {
	cardsInPlay, err := getGameplayCardsForIds(gameplayCardsIdsInPlay)
	if err != nil {
		return []vo.Scores{}, err
	}

	cut, err := getGameplayCardsForIds([]int32{match.Cutgamecardid})
	if err != nil {
		return []vo.Scores{}, err
	}

	for i := 0; i < len(cardsInPlay); i++ {
		details1 := cardDetails(cardsInPlay[i].Value)

		if details1.Value == 11 && cardsInPlay[i].Suit == cut[0].Suit {
			return []vo.Scores{{
				Cards: []vo.GameCard{cardsInPlay[0], cardsInPlay[1]},
				Point: 1,
			}}, nil
		}
	}

	return []vo.Scores{}, nil
}

func getGameplayCardsForIds(ids []int32) ([]vo.GameCard, error) {
	// if len(ids) == 0 {
	// 	return []vo.GameCard{}, nil
	// }

	// // string_ids := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ids)), ","), "[]")
	// cards, err := utils.QueryForCards(ids)
	// if err != nil {
	// 	return []vo.GameCard{}, err
	// }
	// return cards, nil

	return []vo.GameCard{}, nil
}

func scanForThirtyOne(gameplayCardsIdsInPlay []int32) ([]vo.Scores, error) {
	cardsInPlay, err := getGameplayCardsForIds(gameplayCardsIdsInPlay)
	if err != nil {
		return []vo.Scores{}, err
	}

	plays := []vo.GameCard{}
	pointsFound := []vo.Scores{}

	if len(cardsInPlay) < 3 {
		return pointsFound, nil
	}

	details1 := cardDetails(cardsInPlay[0].Value)
	details2 := cardDetails(cardsInPlay[1].Value)
	details3 := cardDetails(cardsInPlay[2].Value)

	if details1.Value+details2.Value+details3.Value == 31 {
		pointsFound = append(pointsFound, vo.Scores{
			Cards: plays,
			Point: 2,
		})
	}

	if len(cardsInPlay) < 4 {
		return pointsFound, nil
	}

	details4 := cardDetails(cardsInPlay[3].Value)

	if details1.Value+details2.Value+details4.Value == 31 {
		pointsFound = append(pointsFound, vo.Scores{
			Cards: plays,
			Point: 2,
		})
	}

	if details1.Value+details3.Value+details4.Value == 31 {
		pointsFound = append(pointsFound, vo.Scores{
			Cards: plays,
			Point: 2,
		})
	}

	if details2.Value+details3.Value+details4.Value == 31 {
		pointsFound = append(pointsFound, vo.Scores{
			Cards: plays,
			Point: 2,
		})
	}

	if details1.Value+details2.Value+details3.Value+details4.Value == 31 {
		pointsFound = append(pointsFound, vo.Scores{
			Cards: plays,
			Point: 2,
		})
	}

	return pointsFound, nil
}

func scanForLastCard(m vo.GameMatch) ([]vo.Scores, error) {
	if len(m.Players) < 2 {
		return []vo.Scores{}, errors.New("too few players")
	}

	playerOneCards, err := getGameplayCardsForIds(m.Players[0].Hand)

	if err != nil {
		return []vo.Scores{}, err
	}

	playerTwoCards, err := getGameplayCardsForIds(m.Players[1].Hand)

	if err != nil {
		return []vo.Scores{}, err
	}

	cardsInPlay, err := getGameplayCardsForIds(cardsInPlay(m.Players))

	if err != nil {
		return []vo.Scores{}, err
	}

	total := 0
	for _, card := range cardsInPlay {
		details := cardDetails(card.Value)
		total = total + details.Value
	}

	playerOneCanPlay := false
	for _, card := range playerOneCards {
		details := cardDetails(card.Value)
		if int(details.Value)+total <= 31 {
			playerOneCanPlay = true
		}
	}

	playerTwoCanPlay := false
	for _, card := range playerTwoCards {
		details := cardDetails(card.Value)
		if int(details.Value)+total <= 31 {
			playerTwoCanPlay = true
		}
	}

	if !playerOneCanPlay && !playerTwoCanPlay {
		return []vo.Scores{
			{
				Point: 1,
			},
		}, nil
	}

	return []vo.Scores{}, nil
}

func scanJackOnCut(match vo.GameMatch) ([]vo.Scores, error) {
	cardsInPlay, err := getGameplayCardsForIds([]int32{match.Cutgamecardid})
	if err != nil || len(cardsInPlay) == 0 {
		return []vo.Scores{}, err
	}

	details := cardDetails(cardsInPlay[0].Value)

	if details.Order == 11 {
		return []vo.Scores{{
			Cards: []vo.GameCard{},
			Point: 2,
		}}, nil
	}

	return []vo.Scores{}, nil
}
