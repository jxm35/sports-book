package bet_selector

import (
	"errors"

	results "github.com/jxm35/go-results"

	"sports-book.com/pkg/config"
	"sports-book.com/pkg/db"
	"sports-book.com/pkg/domain"
)

var ErrInvalidConfig = errors.New("invalid bet placer config provided")

type fixedAmountBetSelector struct {
	minOddsDelta float64
	maxOddsDelta float64
	betAmount    float64
}

func NewFixedAmountBetSelector(minOddsDelta, maxOddsDelta, betAmount float64) BetSelector {
	if betAmount < 0 {
		// return nil, ErrInvalidConfig
		panic(ErrInvalidConfig)
	}
	return &fixedAmountBetSelector{
		minOddsDelta: minOddsDelta,
		maxOddsDelta: maxOddsDelta,
		betAmount:    betAmount,
	}
}

func NewFixedAmountBetSelectorFromConfig() (BetSelector, error) {
	minOddsDelta, found := config.GetConfigVal[float64]("bet_selector.min_odds_delta").Get()
	if !found {
		return nil, ErrInvalidConfig
	}
	maxOddsDelta, found := config.GetConfigVal[float64]("bet_selector.max_odds_delta").Get()
	if !found {
		return nil, ErrInvalidConfig
	}
	betAmount, found := config.GetConfigVal[float64]("bet_selector.bet_amount").Get()
	if !found {
		return nil, ErrInvalidConfig
	}
	if betAmount < 0 {
		return nil, ErrInvalidConfig
	}
	return &fixedAmountBetSelector{
		minOddsDelta: minOddsDelta,
		maxOddsDelta: maxOddsDelta,
		betAmount:    betAmount,
	}, nil
}

func (f *fixedAmountBetSelector) Place1x2Bets(matchId int32, generatedOdds domain.MatchProbability, currentPot float64) results.Option[domain.BetOrder] {
	if f.maxOddsDelta == 0 {
		f.maxOddsDelta = 1
	}
	if f.betAmount > currentPot {
		return results.None[domain.BetOrder]()
	}

	odds := db.GetBestOddsForMatch(matchId)
	bookieImpliedOdds := domain.MatchProbability{
		HomeWin: 1 / odds.HomeWin,
		Draw:    1 / odds.Draw,
		AwayWin: 1 / odds.AwayWin,
	}
	if f.isWithinConstraints(generatedOdds.HomeWin, bookieImpliedOdds.HomeWin) {
		return results.Some(domain.BetOrder{
			MatchId:              matchId,
			Backing:              domain.BackHomeWin,
			BookMaker:            odds.HomeBookie,
			OddsTaken:            odds.HomeWin,
			Amount:               f.betAmount,
			PredictedProbability: generatedOdds.HomeWin,
		})
	}
	if f.isWithinConstraints(generatedOdds.Draw, bookieImpliedOdds.Draw) {
		return results.Some(domain.BetOrder{
			MatchId:              matchId,
			Backing:              domain.BackDraw,
			BookMaker:            odds.DrawBookie,
			OddsTaken:            odds.Draw,
			Amount:               f.betAmount,
			PredictedProbability: generatedOdds.Draw,
		})
	}
	if f.isWithinConstraints(generatedOdds.AwayWin, bookieImpliedOdds.AwayWin) {
		return results.Some(domain.BetOrder{
			MatchId:              matchId,
			Backing:              domain.BackAwayWin,
			BookMaker:            odds.AwayBookie,
			OddsTaken:            odds.AwayWin,
			Amount:               f.betAmount,
			PredictedProbability: generatedOdds.AwayWin,
		})
	}
	return results.None[domain.BetOrder]()
}

func (f *fixedAmountBetSelector) isWithinConstraints(myOdds, bookieOdds float64) bool {
	if myOdds-bookieOdds > f.minOddsDelta && myOdds-bookieOdds < f.maxOddsDelta {
		return true
	}
	return false
}
