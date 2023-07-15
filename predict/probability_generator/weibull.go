package probability_generator

import (
	"math"
	"math/big"
)

type WeibullOddsGenerator struct{}

type alphaArgs struct {
	x float64
	j float64
	c float64
}

var cache = map[alphaArgs]float64{}

const (
	weibullHomeShape = 0.7910484751 //0.9062415
	weibullAwayShape = 0.6016388937 //0.8491849
)

func (p *WeibullOddsGenerator) Generate1x2Probabilities(homeProjected, awayProjected float64) MatchProbability {
	cache = make(map[alphaArgs]float64)
	var homeGoalProb = make(map[int]float64)
	var awayGoalProb = make(map[int]float64)
	for i := 0; i <= 10; i++ {
		homeGoalProb[i] = getGoalProbabilityWeibull(i, homeProjected, weibullHomeShape) // TODO calc variance for home wins
		awayGoalProb[i] = getGoalProbabilityWeibull(i, awayProjected, weibullAwayShape) // TODO cal variance for away wins
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

func getGoalProbabilityWeibull(goals int, projected, shape float64) float64 {
	return weibull(float64(goals), shape, projected)
}

type weibullCountModel struct {
	λ big.Float // lamda rate of 1.5
	c big.Float // shape param 1.56
	t big.Float // time default 1.0
}

//
//func (w *weibullCountModel) pdf(x big.Int) (big.Float, error) {
//	return weibullCountPdf(x, w.λ, w.c)
//}
//
//func weibullCountPdf(x big.Int, λ, c big.Float) (big.Float, error) {
//	t := big.NewFloat(1)
//	p, converged := weibullCountApprox(x, λ, c, t)
//}
//
//func weibullCountApprox(x big.Int, λ, c, t big.Float) {
//	maximum :=
//}

// from python github
func alphas(x, j float64, c float64) float64 {
	if res, ok := cache[alphaArgs{x, j, c}]; ok {
		return res
	}
	if x == 0 {
		res := math.Gamma(c*j+1) / math.Gamma(float64(j)+1)
		cache[alphaArgs{x, j, c}] = res
		return res
	} else if j < x {
		panic("invalid input")
	} else {
		var sum float64
		for m := x - 1; m < j; m++ {
			sum += alphas(x-1, m, c) * math.Gamma(c*j-c*m+1) / math.Gamma(j-m+1)
		}
		cache[alphaArgs{x, j, c}] = sum
		return sum
	}
}

func _weibull(x, c, lambda, t, j float64) float64 {
	return (math.Pow(-1, x+j) * math.Pow(math.Pow(lambda*t, c), j) * alphas(x, j, c)) / math.Gamma(c*j+1)
}

// x = the number we want to predict?, lambda = rate (preset) 1.5, c = shape (preset) 1.56

// weibull returns the probability that there are x goals scored, given lambda predicted (rate of) goals scored
// c denotes the 'shape' of the weibull distribution
func weibull(x, c, lambda float64) float64 {
	var sum float64
	for j := x; j < x+50; j++ {
		sum += _weibull(x, c, lambda, 1, j)
	}
	return sum
}

func cumulativeWeibull(x, c, lambda float64) float64 {
	var sum float64
	for i := 0; i <= int(x); i++ {
		sum += weibull(float64(i), c, lambda)
	}
	return sum
}
