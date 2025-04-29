package core

import (
	"cmp"
	"slices"

	"github.com/bardic/gocrib/queries/queries"

	"github.com/bardic/gocrib/vo"
)

var (
	Zero     int = 0
	One      int = 1
	Two      int = 2
	Three    int = 3
	Four     int = 4
	Five     int = 5
	Six      int = 6
	Seven    int = 7
	Eight    int = 8
	Nine     int = 9
	Ten      int = 10
	Eleven   int = 11
	Twelve   int = 12
	Thirteen int = 13
	Fourteen int = 14
)

func cardsInPlay(players []*vo.GamePlayer) []*int {
	cardIds := []*int{}
	// for _, p := range players {
	// 	cardIds = append(cardIds, p.Hand...)
	// }

	return cardIds
}

func countPegs(scoreVO vo.ScoreMatch) (vo.ScoreResults, error) {
	res := vo.ScoreResults{}

	// c, err := getGameplayCardsForIds(cardsInPlay(m.Players))
	// matchActivePlayerId := 1
	// scoringPlayerId := 1
	// playerHands := []vo.GamePlayer{}

	// if err != nil {
	// 	return vo.ScoreResults{}, err
	// }

	// cut, err := getCut(m)

	// if err != nil {
	// 	return vo.ScoreResults{}, err
	// }

	r, err := scanForThirtyOne(*scoreVO.CardsInPlay)
	if err != nil {
		return vo.ScoreResults{}, err
	}
	res.Results = append(res.Results, r...)

	// if len(m.Turnpasstimestamps) == 0 {
	// 	r, err = scanJackOnCut(cut)
	// 	if err != nil {
	// 		return vo.ScoreResults{}, err
	// 	}
	// 	res.Results = append(res.Results, r...)
	// }

	r, err = scanForFifteens(*scoreVO.CardsInPlay)
	if err != nil {
		return vo.ScoreResults{}, err
	}
	res.Results = append(res.Results, r...)

	// r, err = scanForLastCard(matchActivePlayerId, scoringPlayerId, c, playerHands)
	// if err != nil {
	// 	return vo.ScoreResults{}, err
	// }
	// res.Results = append(res.Results, r...)

	r, err = scanForMatchingKinds(*scoreVO.CardsInPlay)
	if err != nil {
		return vo.ScoreResults{}, err
	}
	res.Results = append(res.Results, r...)

	r, err = scanForRuns(*scoreVO.CardsInPlay)
	if err != nil {
		return vo.ScoreResults{}, err
	}
	res.Results = append(res.Results, r...)

	if scoreVO.ActivePlayerId != nil && scoreVO.Players != nil {
		r, err = scanForLastCard(*scoreVO.ActivePlayerId, *scoreVO.CardsInPlay, *scoreVO.Players)
		if err != nil {
			return vo.ScoreResults{}, err
		}
		res.Results = append(res.Results, r...)
	}

	return vo.ScoreResults{Results: res.Results}, nil
}

func countHand(m vo.GameMatch) (vo.ScoreResults, error) {
	res := vo.ScoreResults{}

	c, err := getGameplayCardsForIds(cardsInPlay(m.Players))
	matchActivePlayerId := 1
	// scoringPlayerId := 1
	playerHands := []vo.GamePlayer{}

	if err != nil {
		return vo.ScoreResults{}, err
	}

	cut, err := getCut(m)

	r, err := scanForThirtyOne(c)
	if err != nil {
		return vo.ScoreResults{}, err
	}
	res.Results = append(res.Results, r...)

	r, err = scanRightJackCut(c, cut)
	if err != nil {
		return vo.ScoreResults{}, err
	}
	res.Results = append(res.Results, r...)

	r, err = scanForFifteens(c)
	if err != nil {
		return vo.ScoreResults{}, err
	}
	res.Results = append(res.Results, r...)

	r, err = scanForLastCard(matchActivePlayerId, c, playerHands)
	if err != nil {
		return vo.ScoreResults{}, err
	}
	res.Results = append(res.Results, r...)

	r, err = scanForMatchingKinds(c)
	if err != nil {
		return vo.ScoreResults{}, err
	}
	res.Results = append(res.Results, r...)

	r, err = scanForRuns(c)
	if err != nil {
		return vo.ScoreResults{}, err
	}
	res.Results = append(res.Results, r...)

	r, err = scanForFlush(c)
	if err != nil {
		return vo.ScoreResults{}, err
	}
	res.Results = append(res.Results, r...)

	return vo.ScoreResults{Results: res.Results}, nil
}

func getCut(m vo.GameMatch) (vo.GameCard, error) {
	return vo.GameCard{

		Card: queries.Card{
			Value: queries.CardvalueAce,
			Suit:  queries.CardsuitClubs,
		},
	}, nil
}

func scanForFlush(cardsInPlay []vo.GameCard) ([]vo.Scores, error) {
	// gameplayCardsInHand, err := getGameplayCardsForIds(cardIdsInHand)
	// if err != nil {
	// 	return []vo.Scores{}, err
	// }

	var flushSuit *queries.Cardsuit
	matchesFlush := true
	for _, card := range cardsInPlay {
		if flushSuit == nil {
			flushSuit = &card.Card.Suit
		} else {
			matchesFlush = *flushSuit == card.Card.Suit
		}
	}

	if matchesFlush {
		return []vo.Scores{
			{
				Point: &Four,
			},
		}, nil
	}

	return []vo.Scores{}, nil
}

func scanForRuns(gameplayCardsInPlay []vo.GameCard) ([]vo.Scores, error) {
	gameplayCardsInPlay = orderHand(gameplayCardsInPlay)

	var pointsFound []vo.Scores

	if len(gameplayCardsInPlay) < 3 {
		return pointsFound, nil
	}

	gameCard1 := gameplayCardsInPlay[0]
	gameCard2 := gameplayCardsInPlay[1]
	gameCard3 := gameplayCardsInPlay[2]

	details1 := cardDetails(gameCard1.Card.Value)
	details2 := cardDetails(gameCard2.Card.Value)
	details3 := cardDetails(gameCard3.Card.Value)

	var pointsToGain *int

	if *details1.Order+1 == *details2.Order && *details2.Order+1 == *details3.Order {
		pointsToGain = &Three
		pointsFound = []vo.Scores{{
			Cards: []vo.GameCard{gameCard1, gameCard2, gameCard3},
			Point: pointsToGain,
		}}

		if len(gameplayCardsInPlay) <= 3 {
			return pointsFound, nil
		}

		gameCard4 := gameplayCardsInPlay[3]
		details4 := cardDetails(gameCard4.Card.Value)

		if *details3.Order+1 == *details4.Order {
			pointsToGain = &Four
			pointsFound = []vo.Scores{{
				Cards: []vo.GameCard{gameCard1, gameCard2, gameCard3, gameCard4},
				Point: pointsToGain,
			}}
		}
	}

	if len(gameplayCardsInPlay) <= 3 {
		return pointsFound, nil
	}

	gameCard4 := gameplayCardsInPlay[3]
	details4 := cardDetails(gameCard4.Card.Value)

	if *details1.Order+1 == *details2.Order &&
		*details2.Order+1 == *details4.Order {

		pointsToGain = &Three
		pointsFound = append(pointsFound, vo.Scores{
			Cards: []vo.GameCard{gameCard1, gameCard2, gameCard4},
			Point: pointsToGain,
		})
	}

	return pointsFound, nil
}

func cardDetails(cardValue queries.Cardvalue) vo.GameCardDetails {
	switch cardValue {
	case queries.CardvalueAce:

		return vo.GameCardDetails{
			Value: &One,
			Order: &One,
		}
	case queries.CardvalueTwo:
		return vo.GameCardDetails{
			Value: &Two,
			Order: &Two,
		}
	case queries.CardvalueThree:
		return vo.GameCardDetails{
			Value: &Three,
			Order: &Three,
		}
	case queries.CardvalueFour:
		return vo.GameCardDetails{
			Value: &Four,
			Order: &Four,
		}
	case queries.CardvalueFive:
		return vo.GameCardDetails{
			Value: &Five,
			Order: &Five,
		}
	case queries.CardvalueSix:
		return vo.GameCardDetails{
			Value: &Six,
			Order: &Six,
		}
	case queries.CardvalueSeven:
		return vo.GameCardDetails{
			Value: &Seven,
			Order: &Seven,
		}
	case queries.CardvalueEight:
		return vo.GameCardDetails{
			Value: &Eight,
			Order: &Eight,
		}
	case queries.CardvalueNine:
		return vo.GameCardDetails{
			Value: &Nine,
			Order: &Nine,
		}
	case queries.CardvalueTen:
		return vo.GameCardDetails{
			Value: &Ten,
			Order: &Ten,
		}
	case queries.CardvalueJack:
		return vo.GameCardDetails{
			Value: &Ten,
			Order: &Eleven,
		}
	case queries.CardvalueQueen:
		return vo.GameCardDetails{
			Value: &Ten,
			Order: &Twelve,
		}
	case queries.CardvalueKing:
		return vo.GameCardDetails{
			Value: &Ten,
			Order: &Thirteen,
		}
	case queries.CardvalueJoker:
		return vo.GameCardDetails{
			Value: &Zero,
			Order: &Zero,
		}
	}

	return vo.GameCardDetails{
		Value: &Zero,
		Order: &Zero,
	}

}

func scanForMatchingKinds(gameplayCardsInPlay []vo.GameCard) ([]vo.Scores, error) {
	gameplayCardsInPlay = orderHand(gameplayCardsInPlay)

	var pointsFound []vo.Scores

	if len(gameplayCardsInPlay) < 2 {
		return pointsFound, nil
	}

	gameCard1 := gameplayCardsInPlay[0]
	gameCard2 := gameplayCardsInPlay[1]

	details1 := cardDetails(gameCard1.Card.Value)
	details2 := cardDetails(gameCard2.Card.Value)

	if details1.Order == details2.Order {
		pointsFound = append(pointsFound, vo.Scores{
			Cards: []vo.GameCard{gameCard1, gameCard2},
			Point: &Two,
		})
	}

	if len(gameplayCardsInPlay) < 3 {
		return pointsFound, nil
	}

	gameCard3 := gameplayCardsInPlay[2]
	details3 := cardDetails(gameCard3.Card.Value)

	if details1.Order == details3.Order {
		pointsFound = append(pointsFound, vo.Scores{
			Cards: []vo.GameCard{gameCard1, gameCard3},
			Point: &Two,
		})
	}

	if details2.Order == details3.Order {
		pointsFound = append(pointsFound, vo.Scores{
			Cards: []vo.GameCard{gameCard2, gameCard3},
			Point: &Two,
		})
	}

	if len(gameplayCardsInPlay) < 4 {
		return pointsFound, nil
	}

	gameCard4 := gameplayCardsInPlay[3]
	details4 := cardDetails(gameCard4.Card.Value)

	if details1.Order == details4.Order {
		pointsFound = append(pointsFound, vo.Scores{
			Cards: []vo.GameCard{gameCard1, gameCard4},
			Point: &Two,
		})
	}

	if details2.Order == details4.Order {
		pointsFound = append(pointsFound, vo.Scores{
			Cards: []vo.GameCard{gameCard2, gameCard4},
			Point: &Two,
		})
	}

	if details3.Order == details4.Order {
		pointsFound = append(pointsFound, vo.Scores{
			Cards: []vo.GameCard{gameCard3, gameCard4},
			Point: &Two,
		})
	}

	return pointsFound, nil
}

func scanForFifteens(cardsInPlay []vo.GameCard) ([]vo.Scores, error) {
	// cardsInPlay, err := getGameplayCardsForIds(gameplayCardsIdsInPlay)
	// if err != nil {
	// 	return []vo.Scores{}, err
	// }

	pointsFound := []vo.Scores{}

	//find if any combination of cardsInPlay equals 15
	for i := 0; i < len(cardsInPlay); i++ {
		for j := i; j < len(cardsInPlay); j++ {
			details1 := cardDetails(cardsInPlay[i].Card.Value)
			details2 := cardDetails(cardsInPlay[j].Card.Value)
			// details3 := cardDetails(cardsInPlay[k].Value)

			if *details1.Value+*details2.Value == 15 {
				pointsFound = append(pointsFound, vo.Scores{
					Cards: []vo.GameCard{cardsInPlay[i], cardsInPlay[j]},
					Point: &Two,
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
				details1 := cardDetails(cardsInPlay[i].Card.Value)
				details2 := cardDetails(cardsInPlay[j].Card.Value)
				details3 := cardDetails(cardsInPlay[k].Card.Value)

				if *details1.Value+*details2.Value+*details3.Value == 15 && i != j && j != k {
					pointsFound = append(pointsFound, vo.Scores{
						Cards: []vo.GameCard{cardsInPlay[i], cardsInPlay[j], cardsInPlay[k]},
						Point: &Two,
					})
				}
			}
		}
	}

	return pointsFound, nil
}

func scanRightJackCut(cardsInPlay []vo.GameCard, cutCard vo.GameCard) ([]vo.Scores, error) {

	// cardsInPlay, err := getGameplayCardsForIds(gameplayCardsIdsInPlay)
	// if err != nil {
	// 	return []vo.Scores{}, err
	// }

	// cut, err := getGameplayCardsForIds([]*int{match.Cutgamecardid})
	// if err != nil {
	// 	return []vo.Scores{}, err
	// }
	for i := 0; i < len(cardsInPlay); i++ {
		details1 := cardDetails(cardsInPlay[i].Card.Value)

		if details1.Order == &Eleven && cardsInPlay[i].Card.Suit == cutCard.Card.Suit {
			return []vo.Scores{{
				Cards: []vo.GameCard{cardsInPlay[0], cardsInPlay[1]},
				Point: &One,
			}}, nil
		}
	}

	return []vo.Scores{}, nil
}

func getGameplayCardsForIds(ids []*int) ([]vo.GameCard, error) {
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

func scanForThirtyOne(cardsInPlay []vo.GameCard) ([]vo.Scores, error) {
	// cardsInPlay, err := getGameplayCardsForIds(gameplayCardsIdsInPlay)
	// if err != nil {
	// 	return []vo.Scores{}, err
	// }

	plays := []vo.GameCard{}
	pointsFound := []vo.Scores{}

	if len(cardsInPlay) < 3 {
		return pointsFound, nil
	}

	details1 := cardDetails(cardsInPlay[0].Card.Value)
	details2 := cardDetails(cardsInPlay[1].Card.Value)
	details3 := cardDetails(cardsInPlay[2].Card.Value)

	if *details1.Value+*details2.Value+*details3.Value == 31 {
		pointsFound = append(pointsFound, vo.Scores{
			Cards: plays,
			Point: &Two,
		})
	}

	if len(cardsInPlay) < 4 {
		return pointsFound, nil
	}

	details4 := cardDetails(cardsInPlay[3].Card.Value)

	if *details1.Value+*details2.Value+*details4.Value == 31 {
		pointsFound = append(pointsFound, vo.Scores{
			Cards: plays,
			Point: &Two,
		})
	}

	if *details1.Value+*details3.Value+*details4.Value == 31 {
		pointsFound = append(pointsFound, vo.Scores{
			Cards: plays,
			Point: &Two,
		})
	}

	if *details2.Value+*details3.Value+*details4.Value == 31 {
		pointsFound = append(pointsFound, vo.Scores{
			Cards: plays,
			Point: &Two,
		})
	}

	if *details1.Value+*details2.Value+*details3.Value+*details4.Value == 31 {
		pointsFound = append(pointsFound, vo.Scores{
			Cards: plays,
			Point: &Two,
		})
	}

	return pointsFound, nil
}

func sum(cards []vo.GameCard) int {
	total := 0
	for _, card := range cards {
		details := cardDetails(card.Card.Value)
		total = total + *details.Value
	}

	return total
}

func scanForLastCard(currentPlayerTurn int, cardsInPlay []vo.GameCard, players []vo.GamePlayer) ([]vo.Scores, error) {

	totalInPlay := sum(cardsInPlay)

	if totalInPlay == 31 {
		return []vo.Scores{
			{
				Point: &Two,
			},
		}, nil
	}

	//sort array of players so that the current player is first
	currnetPlayerIndex := slices.IndexFunc(players, func(i vo.GamePlayer) bool {
		return *i.ID == currentPlayerTurn
	})

	players = append(players[currnetPlayerIndex:], players[:currnetPlayerIndex]...)

	hasPlayableCard := false
	for _, player := range players {
		//if *player.ID == currentPlayerTurn {
		//If if a card in this players hand can be added to the cards in play without xceeding 31
		//then they get a point

		for _, card := range player.Hand {
			details := cardDetails(card.Card.Value)
			if *details.Value+totalInPlay <= 31 {
				hasPlayableCard = true
			}
		}
		//}
	}

	if hasPlayableCard {
		return []vo.Scores{
			{
				Point: &Zero,
			},
		}, nil
	}

	return []vo.Scores{
		{
			Point: &One,
		},
	}, nil
}

func scanForAdditionalPlays(currentPlayerTurn int, cardsInPlay []vo.GameCard, players []vo.GamePlayer) ([]vo.Scores, error) {

	totalInPlay := sum(cardsInPlay)

	if totalInPlay == 31 {
		return []vo.Scores{
			{
				Point: &Two,
			},
		}, nil
	}

	//sort array of players so that the current player is first
	currnetPlayerIndex := slices.IndexFunc(players, func(i vo.GamePlayer) bool {
		return *i.ID == currentPlayerTurn
	})

	players = append(players[currnetPlayerIndex:], players[:currnetPlayerIndex]...)

	for _, player := range players {
		//if *player.ID == currentPlayerTurn {
		//If if a card in this players hand can be added to the cards in play without xceeding 31
		//then they get a point

		for _, card := range player.Hand {
			details := cardDetails(card.Card.Value)
			if *details.Value+totalInPlay <= 31 {
				return []vo.Scores{
					{
						Point: &Zero,
					},
				}, nil
			}
		}
		//}
	}

	return []vo.Scores{
		{
			Point: &One,
		},
	}, nil
}
func scanJackOnCut(cutCard vo.GameCard) ([]vo.Scores, error) {

	if cutCard.Card.Value == queries.CardvalueJack {
		return []vo.Scores{{
			Cards: []vo.GameCard{},
			Point: &One,
		}}, nil

	}

	// details := cardDetails(cutCard.Card.Value)

	// if details.Order == &Eleven {
	// 	return []vo.Scores{{
	// 		Cards: []vo.GameCard{},
	// 		Point: &Two,
	// 	}}, nil
	// }

	return []vo.Scores{}, nil
}

func orderHand(gameplayCardsInPlay []vo.GameCard) []vo.GameCard {
	slices.SortFunc(gameplayCardsInPlay, func(i, j vo.GameCard) int {
		v1 := cardDetails(i.Card.Value).Order
		v2 := cardDetails(j.Card.Value).Order

		return cmp.Compare(*v1, *v2)
	})

	return gameplayCardsInPlay
}
