package probability_generator

import (
	"math"
	"math/big"
)

type PoissonOddsGenerator struct{}

func (p *PoissonOddsGenerator) Generate1x2Probabilities(homeProjected, awayProjected float64) MatchProbability {
	var homeGoalProb = make(map[int]float64)
	var awayGoalProb = make(map[int]float64)
	for i := 0; i <= 10; i++ {
		homeGoalProb[i] = getGoalProbabilityPoisson(i, homeProjected)
		awayGoalProb[i] = getGoalProbabilityPoisson(i, awayProjected)
	}
	matchProb := MatchProbability{
		HomeWin: 0,
		Draw:    0,
		AwayWin: 0,
	}
	for h := 0; h <= 10; h++ {
		for a := 0; a <= 10; a++ { // home win
			if h > a {
				matchProb.HomeWin += homeGoalProb[h] * awayGoalProb[a]
			} else if h == a { // draw
				matchProb.Draw += homeGoalProb[h] * awayGoalProb[a]
			} else { // away win
				matchProb.AwayWin += homeGoalProb[h] * awayGoalProb[a]
			}
		}
	}
	return matchProb
}

func getGoalProbabilityPoisson(desiredGoals int, expected float64) float64 {
	desiredFactorial, _ := new(big.Float).SetInt(factorial(big.NewInt(int64(desiredGoals)))).Float64()
	probDesired := (math.Pow(expected, float64(desiredGoals)) * math.Pow(math.E, -expected)) / desiredFactorial
	return probDesired
}

func factorial(x *big.Int) *big.Int {
	n := big.NewInt(1)
	if x.Cmp(big.NewInt(0)) == 0 {
		return n
	}
	return n.Mul(x, factorial(n.Sub(x, n)))
}
