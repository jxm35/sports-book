package predict

import "sports-book.com/predict/probability_generator"

type Pipeline interface {
	PredictMatch(homeTeam, awayTeam, season int32) (probability_generator.MatchProbability, OddsDelta, error)
}
