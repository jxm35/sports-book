package probability_generator

import (
	"sports-book.com/pkg/domain"
)

type ProbabilityGenerator interface {
	Generate1x2Probabilities(homeProjected, awayProjected float64, league domain.League) domain.MatchProbability
}
