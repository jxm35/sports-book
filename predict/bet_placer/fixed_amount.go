package bet_placer

import (
	"errors"

	results "github.com/jxm35/go-results"

	"sports-book.com/predict/domain"
	"sports-book.com/util"
)

var ErrInvalidConfig = errors.New("invalid bet placer config provided")

type fixedAmountBetPlacer struct {
	minOddsDelta float64
	maxOddsDelta float64
	betAmount    float64
}

func NewFixedAmountBetPlacer(minOddsDelta, maxOddsDelta, betAmount float64) BetPlacer {
	if betAmount < 0 {
		//return nil, ErrInvalidConfig
		panic(ErrInvalidConfig)
	}
	return &fixedAmountBetPlacer{
		minOddsDelta: minOddsDelta,
		maxOddsDelta: maxOddsDelta,
		betAmount:    betAmount,
	}
}

func (f *fixedAmountBetPlacer) Place1x2Bets(matchId int32, generatedOdds domain.MatchProbability, currentPot float64) results.Option[domain.BetOrder] {
	if f.maxOddsDelta == 0 {
		f.maxOddsDelta = 1
	}
	if f.betAmount > currentPot {
		return results.None[domain.BetOrder]()
	}

	odds := util.GetBestOddsForMatch(matchId)
	bookieImpliedOdds := domain.MatchProbability{
		HomeWin: 1 / odds.HomeWin,
		Draw:    1 / odds.Draw,
		AwayWin: 1 / odds.AwayWin,
	}
	if f.isWithinConstraints(generatedOdds.HomeWin, bookieImpliedOdds.HomeWin) {
		return results.Some(domain.BetOrder{
			MatchId:   matchId,
			Backing:   domain.BackHomeWin,
			BookMaker: odds.HomeBookie,
			OddsTaken: odds.HomeWin,
			Amount:    f.betAmount,
		})
	}
	if f.isWithinConstraints(generatedOdds.Draw, bookieImpliedOdds.Draw) {
		return results.Some(domain.BetOrder{
			MatchId:   matchId,
			Backing:   domain.BackDraw,
			BookMaker: odds.DrawBookie,
			OddsTaken: odds.Draw,
			Amount:    f.betAmount,
		})
	}
	if f.isWithinConstraints(generatedOdds.AwayWin, bookieImpliedOdds.AwayWin) {
		return results.Some(domain.BetOrder{
			MatchId:   matchId,
			Backing:   domain.BackAwayWin,
			BookMaker: odds.AwayBookie,
			OddsTaken: odds.AwayWin,
			Amount:    f.betAmount,
		})
	}
	return results.None[domain.BetOrder]()
}

func (f *fixedAmountBetPlacer) isWithinConstraints(myOdds, bookieOdds float64) bool {
	if myOdds-bookieOdds > f.minOddsDelta && myOdds-bookieOdds < f.maxOddsDelta {
		return true
	}
	return false
}
