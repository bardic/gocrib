package route

import (
	"fmt"
	"testing"

	"github.com/bardic/cribbage/server/model"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestHelloName(t *testing.T) {
	hand := []int{10, 10, 11, 11}

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

func TestPair(t *testing.T) {
	hand := []int{9, 9, 11, 13}

	scores, err := scanForMatchingKinds(hand, []model.Scores{})
	fmt.Println(scores[0])
	if scores[0].Point != 2 || err != nil {
		t.Fatalf(`scanForThirtyOne(hand) = %q, %v, want match for %#q, nil`, scores, err, 2)
	}
}
