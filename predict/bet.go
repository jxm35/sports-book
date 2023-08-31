package predict

type BackingType string

var BackHomeWin BackingType = "BackHomeWin"
var BackDraw BackingType = "BackDraw"
var BackAwayWin BackingType = "BackAwayWin"

type BetOrder struct {
	MatchId   int32
	Backing   BackingType
	BookMaker string
	OddsTaken float64
	Amount    float64
}

func HandleOddsDelta(delta OddsDelta, matchId int32) *BetOrder {
	if delta.HomeWinDelta > 0.5 {
		return &BetOrder{
			MatchId:   matchId,
			Backing:   BackHomeWin,
			BookMaker: delta.HomeBookie,
			OddsTaken: delta.HomeWinOdds,
			Amount:    1,
		}
	}
	if delta.DrawDelta > 0.5 {
		return &BetOrder{
			MatchId:   matchId,
			Backing:   BackDraw,
			BookMaker: delta.DrawBookie,
			OddsTaken: delta.DrawOdds,
			Amount:    1,
		}
	}
	if delta.AwayWinDelta > 0.5 {
		return &BetOrder{
			MatchId:   matchId,
			Backing:   BackAwayWin,
			BookMaker: delta.AwayBookie,
			OddsTaken: delta.AwayWinOdds,
			Amount:    1,
		}
	}
	return nil
}

func kellyCriterion(probability, decimalOdds float64) float64 {
	return (probability / 1) - ((1 - probability) / (decimalOdds - 1)) // odds of 1.5 will return a 50% profit hence 0.5
}
