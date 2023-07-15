package probability_generator

type ProbabilityGenerator interface {
	Generate1x2Probabilities(homeProjected, awayProjected float64) MatchProbability
}

type MatchProbability struct {
	HomeWin float64
	Draw    float64
	AwayWin float64
}
