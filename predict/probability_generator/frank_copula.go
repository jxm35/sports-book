package probability_generator

import (
	"math"
	"sports-book.com/predict/domain"
)

type FrankCopulaOddsGenerator struct {
	KValue float64 // between -1, 1 -- -0.498181645745231 gives -2138.613820304843 (with weibull) -0.24999990725792537 gives -2121.519384258721 (poisson)
}

//func FindFranksValue() {
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
//	var matches = make([]model.Match, 0)
//	for i := 2017; i <= 2022; i++ {
//		yearMatches, err := util.GetMatchesInSeason(int32(i))
//		if err != nil {
//			panic(err)
//		}
//		matches = append(matches, yearMatches...)
//	}
//	maxFunc := func(match model.Match, k float64) float64 {
//		l := &goals_predictor.LastSeasonXgGoalPredictor{}
//		homeExp, awayExp, err := l.PredictScore(match.HomeTeam, match.AwayTeam, competitionYearMap[match.Competition])
//		if errors.Is(err, goals_predictor.ErrNoPreviousData) {
//			return 0
//		}
//		probs := frankWeibullLikelihood(
//			int(match.HomeGoals), int(match.AwayGoals),
//			homeExp, awayExp,
//			1, 1,
//			k,
//		)
//		return math.Log10(probs)
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

func (f *FrankCopulaOddsGenerator) Generate1x2Probabilities(homeProjected, awayProjected float64) domain.MatchProbability {
	if len(cache) == 0 {
		cache = make(map[alphaArgs]float64)
	}
	matchProb := domain.MatchProbability{
		HomeWin: 0,
		Draw:    0,
		AwayWin: 0,
	}
	for h := 0; h <= 10; h++ {
		for a := 0; a <= 10; a++ {
			prob := frankWeibullLikelihood(h, a, homeProjected, awayProjected, weibullHomeShape, weibullAwayShape, f.KValue)
			if h > a {
				matchProb.HomeWin += prob
			} else if h == a {
				matchProb.Draw += prob
			} else {
				matchProb.AwayWin += prob
			}
		}
	}
	return matchProb
}

func frankCopula(u, v, k float64) float64 {
	return ((-1) / k) * math.Log(1+(((math.Exp(-k*u)-1)*(math.Exp(-k*v)-1))/(math.Exp(-k)-1)))
}

func frankWeibullLikelihood(homeGoals, awayGoals int, expHome, expAway float64, cHome, cAway float64, k float64) float64 {
	if k == 0 {
		k = 0.00000000001
	}
	x1 := cumulativeWeibull(float64(homeGoals), cHome, expHome)
	x2 := cumulativeWeibull(float64(awayGoals), cAway, expAway)
	x3 := cumulativeWeibull(float64(homeGoals-1), cHome, expHome)
	x4 := cumulativeWeibull(float64(awayGoals-1), cAway, expAway)

	return frankCopula(x1, x2, k) - frankCopula(x3, x2, k) - frankCopula(x1, x4, k) + frankCopula(x3, x4, k)
}
