package predict

import (
	"errors"
	"fmt"
	"sports-book.com/predict/goals_predictor"
	"sports-book.com/predict/probability_generator"
	"sports-book.com/util"
)

var ErrTeamNotFound = errors.New("team not found")

type pipelineImpl struct {
	predictor            goals_predictor.GoalsPredictor
	probabilityGenerator probability_generator.ProbabilityGenerator
}

type Odds1x2 struct {
	HomeWin float64
	Draw    float64
	AwayWin float64
}

func (p *pipelineImpl) PredictMatch(homeTeam, awayTeam, season int32) (probability_generator.MatchProbability, OddsDelta, error) {
	homeGoalsPredicted, awayGoalsPredicted, err := p.predictor.PredictScore(homeTeam, awayTeam, season)
	if err != nil {
		return probability_generator.MatchProbability{}, OddsDelta{}, err
	}
	matchProbabilities := p.probabilityGenerator.Generate1x2Probabilities(homeGoalsPredicted, awayGoalsPredicted)
	fmt.Printf("my probabilities: %+v\n", matchProbabilities)

	bestOdds := util.GetBestOdds(homeTeam, awayTeam, season)
	fmt.Printf("best odds:%+v\n", bestOdds)

	impliedOdds := Odds1x2{
		HomeWin: 1 / matchProbabilities.HomeWin,
		Draw:    1 / matchProbabilities.Draw,
		AwayWin: 1 / matchProbabilities.AwayWin,
	}

	delta := getOddsDelta(impliedOdds, bestOdds)
	fmt.Printf("odds delta:%+v\n", delta)

	return matchProbabilities, delta, nil
}
