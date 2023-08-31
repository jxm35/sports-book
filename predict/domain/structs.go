package domain

type BackingType string

var BackHomeWin BackingType = "BackHomeWin"
var BackDraw BackingType = "BackDraw"
var BackAwayWin BackingType = "BackAwayWin"

type MatchProbability struct {
	HomeWin float64
	Draw    float64
	AwayWin float64
}

type BetOrder struct {
	MatchId   int32
	Backing   BackingType
	BookMaker string
	OddsTaken float64
	Amount    float64
}
