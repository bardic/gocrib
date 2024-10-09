package state

import (
	"github.com/bardic/gocrib/model"
	"github.com/charmbracelet/bubbles/textinput"
)

var AccountId int
var ActiveMatchId int
var ActiveMatch model.GameMatch
var ActiveDeck *model.GameDeck
var CurrentAction model.GameAction
var CurrentHandModifier model.HandModifier

var CutInput textinput.Model
