package bet_placer

import (
	results "github.com/jxm35/go-results"

	"sports-book.com/predict/domain"
	"sports-book.com/util"
)

type kellyCriterionBetPlacer struct {
	maxPercentBet float64
	minOddsDelta  float64
	maxOddsDelta  float64
	amountFunc    func(maxPercentBet, kellysProportion float64) float64
}

func NewKellyCriterionBetPlacer(minOddsDelta, maxOddsDelta, maxPercentBet float64, linearAmounts bool) BetPlacer {
	if maxOddsDelta == 0 {
		maxOddsDelta = 1
	}
	if maxPercentBet == 0 {
		maxPercentBet = 1
	}
	amountFunc := min
	if linearAmounts {
		amountFunc = func(maxPercentBet, kellysProportion float64) float64 {
			return maxPercentBet * kellysProportion
		}
	}
	return &kellyCriterionBetPlacer{
		minOddsDelta:  minOddsDelta,
		maxOddsDelta:  maxOddsDelta,
		maxPercentBet: maxPercentBet,
		amountFunc:    amountFunc,
	}
}

func (k *kellyCriterionBetPlacer) Place1x2Bets(matchId int32, generatedOdds domain.MatchProbability, currentPot float64) results.Option[domain.BetOrder] {

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
			Amount:    k.amountFunc(k.maxPercentBet, kellyCriterion(generatedOdds.HomeWin, odds.HomeWin)) * currentPot,
		})
	}
	if k.isWithinConstraints(generatedOdds.Draw, bookieImpliedOdds.Draw) {
		return results.Some(domain.BetOrder{
			MatchId:   matchId,
			Backing:   domain.BackDraw,
			BookMaker: odds.DrawBookie,
			OddsTaken: odds.Draw,
			Amount:    k.amountFunc(k.maxPercentBet, kellyCriterion(generatedOdds.Draw, odds.Draw)) * currentPot,
		})
	}
	if k.isWithinConstraints(generatedOdds.AwayWin, bookieImpliedOdds.AwayWin) {
		return results.Some(domain.BetOrder{
			MatchId:   matchId,
			Backing:   domain.BackAwayWin,
			BookMaker: odds.AwayBookie,
			OddsTaken: odds.AwayWin,
			Amount:    k.amountFunc(k.maxPercentBet, kellyCriterion(generatedOdds.AwayWin, odds.AwayWin)) * currentPot,
		})
	}
	return results.None[domain.BetOrder]()
}

// isWithinConstraints returns true if the difference between the supplied odds and the bookmakers odds is within the
// constraints dictated by the kelly criterion bet placer.
func (k *kellyCriterionBetPlacer) isWithinConstraints(myOdds, bookieOdds float64) bool {
	if myOdds-bookieOdds > k.minOddsDelta && myOdds-bookieOdds < k.maxOddsDelta {
		return true
	}
	return false
}

// kellyCriterion takes in probability and decimal odds and returns a proportion of the current pot to bet
func kellyCriterion(probability, decimalOdds float64) float64 {
	return (probability / 1) - ((1 - probability) / (decimalOdds - 1)) // odds of 1.5 will return a 50% profit hence 0.5
}

// min returns the lowest out of the two provided floats
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
