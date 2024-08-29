package route

import (
	"fmt"
	"testing"

	"github.com/bardic/cribbage/server/model"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestThirtyOne(t *testing.T) {
	hand := []int{11, 12, 24, 25}

	scores, err := scanForThirtyOne(hand)
	if len(scores) != 2 || err != nil {
		t.Fatalf(`scanForThirtyOne(hand) = %q, %v, want match for %#q, nil`, scores, err, 2)
	}
}

func TestRunOfFour(t *testing.T) { //error: found two runs. one contains the other
	hand := []int{9, 10, 11, 12}

	scores, err := scanForRuns(hand)
	if scores[0].Point != 4 || err != nil {
		t.Fatalf(`scanForThirtyOne(hand) = %q, %v, want match for %#q, nil`, scores, err, 2)
	}
}

func TestRunOfThree(t *testing.T) {
	hand := []int{9, 10, 11, 13}

	scores, err := scanForRuns(hand)
	if scores[0].Point != 3 || err != nil {
		t.Fatalf(`scanForThirtyOne(hand) = %q, %v, want match for %#q, nil`, scores, err, 2)
	}
}

func TestTwoRunsOfThree(t *testing.T) {
	hand := []int{9, 10, 11, 24}

	scores, err := scanForRuns(hand)
	if len(scores) != 2 || err != nil {
		t.Fatalf(`scanForThirtyOne(hand) = %q, %v, want match for %#q, nil`, scores, err, 2)
	}

	if scores[0].Point != 3 || err != nil {
		t.Fatalf(`scanForThirtyOne(hand) = %q, %v, want match for %#q, nil`, scores, err, 2)
	}

	if scores[1].Point != 3 || err != nil {
		t.Fatalf(`scanForThirtyOne(hand) = %q, %v, want match for %#q, nil`, scores, err, 2)
	}
}

func TestPair(t *testing.T) {
	hand := []int{9, 10, 11, 24}

	scores, err := scanForMatchingKinds(hand)
	fmt.Println(scores[0])
	if scores[0].Point != 2 || err != nil {
		t.Fatalf(`scanForThirtyOne(hand) = %q, %v, want match for %#q, nil`, scores, err, 2)
	}
}

func TestPairThree(t *testing.T) {
	hand := []int{9, 37, 11, 24}

	scores, err := scanForMatchingKinds(hand)
	if len(scores) != 3 || err != nil {
		t.Fatalf(`scanForThirtyOne(hand) = %q, %v, want match for %#q, nil`, scores, err, 2)
	}
}

func TestFlush(t *testing.T) {

}

func TestOneFifthteens(t *testing.T) {
	hand := []int{8, 9, 11, 10}

	scores, err := scanForFifthteens(hand)
	if len(scores) != 1 || err != nil {
		t.Fatalf(`scanForThirtyOne(hand) = %q, %v, want match for %#q, nil`, scores, err, 2)
	}
}

func TestTwoFifthteens(t *testing.T) {
	hand := []int{8, 9, 22, 19}

	scores, err := scanForFifthteens(hand)
	if len(scores) != 2 || err != nil {
		t.Fatalf(`scanForThirtyOne(hand) = %q, %v, want match for %#q, nil`, scores, err, 2)
	}
}

func TestThreeFifthteens(t *testing.T) {
	hand := []int{8, 9, 22, 35}

	scores, err := scanForFifthteens(hand)
	if len(scores) != 3 || err != nil {
		t.Fatalf(`scanForThirtyOne(hand) = %q, %v, want match for %#q, nil`, scores, err, 2)
	}
}

func TestFourFifthteens(t *testing.T) {
	hand := []int{8, 9, 21, 22}

	scores, err := scanForFifthteens(hand)
	if len(scores) != 4 || err != nil {
		t.Fatalf(`scanForThirtyOne(hand) = %q, %v, want match for %#q, nil`, scores, err, 2)
	}
}

func TestTwoFifthteensOfThree(t *testing.T) {
	hand := []int{8, 21, 2, 15}

	scores, err := scanForFifthteens(hand)
	if len(scores) != 2 || err != nil {
		t.Fatalf(`scanForThirtyOne(hand) = %q, %v, want match for %#q, nil`, scores, err, 2)
	}
}

func TestThreeFifthteensOfMix(t *testing.T) {
	hand := []int{8, 9, 2, 21}

	scores, err := scanForFifthteens(hand)
	if len(scores) != 3 || err != nil {
		t.Fatalf(`scanForThirtyOne(hand) = %q, %v, want match for %#q, nil`, scores, err, 2)
	}
}

func TestRightJack(t *testing.T) {
	hand := []int{8, 9, 2, 12}
	match := model.Match{
		CutGameCardId: 1,
	}

	scores, err := scanRightJackCut(hand, match)
	if len(scores) != 1 || err != nil {
		t.Fatalf(`scanForThirtyOne(hand) = %q, %v, want match for %#q, nil`, scores, err, 2)
	}
}

func TestJackOnCut(t *testing.T) {
	match := model.Match{
		CutGameCardId: 12,
	}

	scores, err := scanJackOnCut(match)
	if len(scores) != 1 || err != nil {
		t.Fatalf(`scanForThirtyOne(hand) = %q, %v, want match for %#q, nil`, scores, err, 2)
	}
}
