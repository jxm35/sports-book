package probability_generator

import (
	"errors"
	"fmt"
	"math"
	"math/big"

	"gonum.org/v1/gonum/diff/fd"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/optimize"

	"sports-book.com/pkg/model"
	"sports-book.com/pkg/predict/domain"
	goals_predictor2 "sports-book.com/pkg/predict/goals_predictor"
	"sports-book.com/pkg/util"
)

type WeibullOddsGenerator struct{}

type alphaArgs struct {
	x float64
	j float64
	c float64
}

var cache = map[alphaArgs]float64{}

const (
	// 0.9213043557296844 predicted by maximising with gradient threshold
	weibullHomeShape = 0.6963436844 // 0.7910484751 //0.9062415
	// 0.888478275810224 predicted by maximising with gradient threshold
	weibullAwayShape = 0.6963436844 // 0.6016388937 //0.8491849
)

var weibullShape = map[string]float64{
	domain.LeagueEPL:        0.6963436844,
	domain.LeagueLaLiga:     0.7198451013,
	domain.LeagueSerieA:     0.6867467538,
	domain.LeagueBundesliga: 0.7872960482,
	"general":               0.7225578969,
}

func FindWeibullShapes() {
	league := "epl"
	competitionYearMap := map[int32]int32{
		1:  2018,
		2:  2019,
		3:  2020,
		4:  2021,
		5:  2022,
		7:  2014,
		8:  2015,
		9:  2016,
		10: 2017,
	}
	cache = make(map[alphaArgs]float64)
	matches := make([]model.Match, 0)
	for i := 2017; i <= 2022; i++ {
		yearMatches, err := util.GetMatchesInSeason(int32(i))
		if err != nil {
			panic(err)
		}
		matches = append(matches, yearMatches...)
	}
	maxFunc := func(match model.Match, shape float64) float64 {
		l := &goals_predictor2.LastSeasonXgGoalPredictor{}
		_, awayExp, err := l.PredictScore(match.HomeTeam, match.AwayTeam, competitionYearMap[match.Competition], league, match.Date, match.ID)
		if errors.Is(err, goals_predictor2.ErrNoPreviousData) {
			return 0
		}
		probs := getGoalProbabilityWeibull(int(match.AwayGoals), awayExp, shape)
		return math.Log10(probs)
	}
	minFunc := func(x []float64) float64 {
		var sum float64
		for _, match := range matches {
			val := maxFunc(match, x[0])
			if !math.IsNaN(val) {
				sum += val
			}
		}
		return -1 * sum
	}
	p := optimize.Problem{
		Func: minFunc,
		Grad: func(grad, x []float64) {
			fd.Gradient(grad, minFunc, x, nil)
		},
		Hess: func(hess *mat.SymDense, x []float64) {
			fd.Hessian(hess, minFunc, x, nil)
		},
		Status: nil,
	}
	res, err := optimize.Minimize(p, []float64{1}, nil, nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("------- result found -------")
	fmt.Println(res)
}

func (p *WeibullOddsGenerator) Generate1x2Probabilities(homeProjected, awayProjected float64, league string) domain.MatchProbability {
	cache = make(map[alphaArgs]float64)
	homeGoalProb := make(map[int]float64)
	awayGoalProb := make(map[int]float64)
	// homeShape := util.GetHomexGVariance(season - 1)
	// awayShape := util.GetAwayxGVariance(season - 1)
	shape, ok := weibullShape[league]
	if !ok {
		panic("invalid league")
	}
	for i := 0; i <= 10; i++ {
		homeGoalProb[i] = getGoalProbabilityWeibull(i, homeProjected, shape) // TODO calc variance for home wins
		awayGoalProb[i] = getGoalProbabilityWeibull(i, awayProjected, shape) // TODO cal variance for away wins
	}
	matchProb := domain.MatchProbability{
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
