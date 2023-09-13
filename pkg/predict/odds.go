package predict

import (
	"sports-book.com/pkg/util"
)

type OddsDelta struct {
	HomeBookie   string
	HomeWinDelta float64
	HomeWinOdds  float64

	DrawBookie string
	DrawDelta  float64
	DrawOdds   float64

	AwayBookie   string
	AwayWinDelta float64
	AwayWinOdds  float64
}

func getOddsDelta(impliedOdds Odds1x2, bookiesOdds util.OddsChecka) OddsDelta {
	delta := OddsDelta{
		HomeBookie:   bookiesOdds.HomeBookie,
		HomeWinDelta: bookiesOdds.HomeWin - impliedOdds.HomeWin,
		HomeWinOdds:  bookiesOdds.HomeWin,

		DrawBookie: bookiesOdds.DrawBookie,
		DrawDelta:  bookiesOdds.Draw - impliedOdds.Draw,
		DrawOdds:   bookiesOdds.Draw,

		AwayBookie:   bookiesOdds.AwayBookie,
		AwayWinDelta: bookiesOdds.AwayWin - impliedOdds.AwayWin,
		AwayWinOdds:  bookiesOdds.AwayWin,
	}
	return delta
}
