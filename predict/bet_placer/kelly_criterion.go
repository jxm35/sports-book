package bet_placer

import (
	"github.com/jxm35/go-results"

	"sports-book.com/predict/domain"
	"sports-book.com/util"
)

type KellyCriterionBetPlacer struct {
	MaxPercentBetted float64
	MinOddsDelta     float64
	MaxOddsDelta     float64
	LinearAmounts    bool
}

func (k *KellyCriterionBetPlacer) isWithinConstraints(myOdds, bookieOdds float64) bool {
	if myOdds-bookieOdds > k.MinOddsDelta && myOdds-bookieOdds < k.MaxOddsDelta {
		return true
	}
	return false
}

func (k *KellyCriterionBetPlacer) Place1x2Bets(matchId int32, generatedOdds domain.MatchProbability, currentPot float64) results.Option[domain.BetOrder] {
	if k.MaxOddsDelta == 0 {
		k.MaxOddsDelta = 1
	}
	if k.MaxPercentBetted == 0 {
		k.MaxPercentBetted = 1
	}
	amountFunc := min
	if k.LinearAmounts {
		amountFunc = func(a, b float64) float64 {
			return a * b
		}
	}

	odds := util.GetBestOddsForMatch(matchId)
	bookieImpliedOdds := domain.MatchProbability{
		HomeWin: 1 / odds.HomeWin,
		Draw:    1 / odds.Draw,
		AwayWin: 1 / odds.AwayWin,
	}
	if k.isWithinConstraints(generatedOdds.HomeWin, bookieImpliedOdds.HomeWin) {
		return results.Some(domain.BetOrder{
			MatchId:   matchId,
			Backing:   domain.BackHomeWin,
			BookMaker: odds.HomeBookie,
			OddsTaken: odds.HomeWin,
			Amount:    amountFunc(k.MaxPercentBetted, kellyCriterion(generatedOdds.HomeWin, odds.HomeWin, currentPot)),
		})
	}
	if k.isWithinConstraints(generatedOdds.Draw, bookieImpliedOdds.Draw) {
		return results.Some(domain.BetOrder{
			MatchId:   matchId,
			Backing:   domain.BackDraw,
			BookMaker: odds.DrawBookie,
			OddsTaken: odds.Draw,
			Amount:    amountFunc(k.MaxPercentBetted, kellyCriterion(generatedOdds.Draw, odds.Draw, currentPot)),
		})
	}
	if k.isWithinConstraints(generatedOdds.AwayWin, bookieImpliedOdds.AwayWin) {
		return results.Some(domain.BetOrder{
			MatchId:   matchId,
			Backing:   domain.BackAwayWin,
			BookMaker: odds.AwayBookie,
			OddsTaken: odds.AwayWin,
			Amount:    amountFunc(k.MaxPercentBetted, kellyCriterion(generatedOdds.AwayWin, odds.AwayWin, currentPot)),
		})
	}
	return results.None[domain.BetOrder]()
}

func kellyCriterion(probability, decimalOdds, currentPot float64) float64 {
	return currentPot * ((probability / 1) - ((1 - probability) / (decimalOdds - 1))) // odds of 1.5 will return a 50% profit hence 0.5
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
