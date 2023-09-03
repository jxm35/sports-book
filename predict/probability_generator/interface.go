package probability_generator

import "sports-book.com/predict/domain"

type ProbabilityGenerator interface {
	Generate1x2Probabilities(homeProjected, awayProjected float64, league string) domain.MatchProbability
}
