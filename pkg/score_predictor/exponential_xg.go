package score_predictor

import (
	"fmt"
	"math"
	"time"

	"sports-book.com/pkg/domain"
	"sports-book.com/pkg/entity"
)

type ExponentialXgGoalPredictor struct {
	MinMatches  int
	DecayFactor float64 // 0.004595630985427014 --
	// -0.0011613404647386283 for poisson x copula @  -0.2499999073
	// -0.000401299106310568 for weibull x copula @ -0.4981816457
}

//
//func FindDecayFactor() {
//	competitionYearMap := map[int32]int32{
//		1:  2018,
//		2:  2019,
//		3:  2020,
//		4:  2021,
//		5:  2022,
//		7:  2014,
//		8:  2015,
//		9:  2016,
//		10: 2017,
//	}
//	gen := &probability_generator.FrankCopulaOddsGenerator{KValue: -0.4981816457}
//	var matches = make([]model.Match, 0)
//	for i := 2017; i <= 2022; i++ {
//		yearMatches, err := entity.GetMatchesInSeason(int32(i))
//		if err != nil {
//			panic(err)
//		}
//		matches = append(matches, yearMatches...)
//	}
//	maxFunc := func(match model.Match, decay float64) float64 {
//		e := &ExponentialXgGoalPredictor{
//			MinMatches:  19,
//			DecayFactor: decay,
//		}
//		// get the probability of a home/draw/away win
//		homeExpected, awayExpected, _ := e.PredictScore(match.HomeTeam, match.AwayTeam, competitionYearMap[match.Competition])
//		probs := gen.Generate1x2Probabilities(homeExpected, awayExpected)
//		// return log10 of the correct result
//		if match.HomeGoals > match.AwayGoals {
//			return math.Log10(probs.HomeWin)
//		} else if match.HomeGoals == match.AwayGoals {
//			return math.Log10(probs.Draw)
//		} else {
//			return math.Log10(probs.AwayWin)
//		}
//	}
//	minFunc := func(x []float64) float64 {
//		var sum float64
//		for _, match := range matches {
//			val := maxFunc(match, x[0])
//			if !math.IsNaN(val) {
//				sum += val
//			}
//		}
//		return -1 * sum
//	}
//	p := optimize.Problem{
//		Func: minFunc,
//		Grad: func(grad, x []float64) {
//			fd.Gradient(grad, minFunc, x, nil)
//		},
//		Hess: func(hess *mat.SymDense, x []float64) {
//			fd.Hessian(hess, minFunc, x, nil)
//		},
//		Status: nil,
//	}
//	res, err := optimize.Minimize(p, []float64{0}, nil, nil)
//	if err != nil {
//		fmt.Println(err)
//	}
//	fmt.Println(res)
//}

func (e *ExponentialXgGoalPredictor) getWeight(currentDate, matchDate time.Time) float64 {
	dayDiff := currentDate.Sub(matchDate).Hours() / 24
	halfWeekDiff := dayDiff / 3.5
	return math.Exp(-e.DecayFactor * halfWeekDiff)
}

func (e *ExponentialXgGoalPredictor) PredictScore(homeTeam, awayTeam, season int32, league domain.League, date time.Time) (float64, float64, error) {
	if e.MinMatches <= 0 {
		panic("invalid minimum matches to test")
	}

	// calculate standard for the year before
	seasonStats, err := entity.GetSeasonDetails(season-1, league)
	if err != nil {
		return -1, -1, err
	}
	if seasonStats.MatchCount == 0 || seasonStats.MatchCount%2 != 0 {
		return -1, -1, ErrInvalidSeason
	}
	avgHomeXg := seasonStats.TotalHomexG / float64(seasonStats.MatchCount)
	avgAwayXg := seasonStats.TotalAwayxG / float64(seasonStats.MatchCount)
	avgHomeGoalsConceded := avgAwayXg
	avgAwayGoalsConceded := avgHomeXg

	// calculate home team's strengths
	homeMatches, err := entity.GetTeamHomeMatchesSince(homeTeam, time.Now().AddDate(-3, 0, 0))
	if err != nil {
		return -1, -1, err
	}
	if len(homeMatches) < e.MinMatches {
		return -1, -1, ErrNoPreviousData
	}

	var sumScoredAtHome, sumConcededAtHome, sumWeightsHome float64
	for _, match := range homeMatches {
		weight := e.getWeight(time.Now(), match.Date)
		sumScoredAtHome += match.HomeExpectedGoals * weight
		sumConcededAtHome += match.AwayExpectedGoals * weight
		sumWeightsHome += weight
	}
	homeWeightedxG := sumScoredAtHome / sumWeightsHome
	homeWeightedxGConceded := sumConcededAtHome / sumWeightsHome

	// calculate away team's strengths
	awayMatches, err := entity.GetTeamAwayMatchesSince(awayTeam, time.Now().AddDate(-3, 0, 0))
	if err != nil {
		return -1, -1, err
	}
	if len(awayMatches) < e.MinMatches {
		return -1, -1, ErrNoPreviousData
	}

	var sumScoredAway, sumConcededAway, sumWeightsAway float64
	for _, match := range homeMatches {
		weight := e.getWeight(time.Now(), match.Date)
		sumScoredAway += match.AwayExpectedGoals * weight
		sumConcededAway += match.HomeExpectedGoals * weight
		sumWeightsAway += weight
	}
	awayWeightedxG := sumScoredAway / sumWeightsAway
	awayWeightedxGConceded := sumConcededAway / sumWeightsAway

	homeAttackStrength := homeWeightedxG / avgHomeXg
	homeDefenseStrength := homeWeightedxGConceded / avgHomeGoalsConceded

	awayAttackStrength := awayWeightedxG / avgAwayXg
	awayDefenseStrength := awayWeightedxGConceded / avgAwayGoalsConceded

	// use strengths to project home and away goals
	projectedHomeGoals := homeAttackStrength * awayDefenseStrength * avgHomeXg
	projectedAwayGoals := awayAttackStrength * homeDefenseStrength * avgAwayXg

	fmt.Printf("%d: %f | %d: %f", homeTeam, projectedHomeGoals, awayTeam, projectedAwayGoals)
	return projectedHomeGoals, projectedAwayGoals, nil
}
