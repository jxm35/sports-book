package probability_generator

import (
	"math"
)

type FrankCopulaOddsGenerator struct{}

const franksK = 1

func (f *FrankCopulaOddsGenerator) Generate1x2Probabilities(homeProjected, awayProjected float64) MatchProbability {
	cache = make(map[alphaArgs]float64)
	matchProb := MatchProbability{
		HomeWin: 0,
		Draw:    0,
		AwayWin: 0,
	}
	for h := 0; h <= 10; h++ {
		for a := 0; a <= 10; a++ {
			prob := frankWeibullLikelihood(h, a, homeProjected, awayProjected, weibullHomeShape, weibullAwayShape)
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

func frankWeibullLikelihood(homeGoals, awayGoals int, expHome, expAway float64, cHome, cAway float64) float64 {
	x1 := cumulativeWeibull(float64(homeGoals), cHome, expHome)
	x2 := cumulativeWeibull(float64(awayGoals), cAway, expAway)
	x3 := cumulativeWeibull(float64(homeGoals-1), cHome, expHome)
	x4 := cumulativeWeibull(float64(awayGoals-1), cAway, expAway)

	return frankCopula(x1, x2, franksK) - frankCopula(x3, x2, franksK) - frankCopula(x1, x4, franksK) + frankCopula(x3, x4, franksK)
}
