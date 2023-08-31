package predict

import (
	"errors"
	"fmt"

	"github.com/jxm35/go-results"

	"sports-book.com/predict/bet_placer"
	"sports-book.com/predict/domain"
	"sports-book.com/predict/goals_predictor"
	"sports-book.com/predict/probability_generator"
	"sports-book.com/util"
)

var ErrTeamNotFound = errors.New("team not found")

type pipelineImpl struct {
	predictor            goals_predictor.GoalsPredictor
	probabilityGenerator probability_generator.ProbabilityGenerator
	betPlacer            bet_placer.BetPlacer
}

type Odds1x2 struct {
	HomeWin float64
	Draw    float64
	AwayWin float64
}

func (p *pipelineImpl) PredictMatch(homeTeam, awayTeam, season int32, league string) (domain.MatchProbability, OddsDelta, error) {
	homeGoalsPredicted, awayGoalsPredicted, err := p.predictor.PredictScore(homeTeam, awayTeam, season, league)
	if err != nil {
		return domain.MatchProbability{}, OddsDelta{}, err
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

func (p *pipelineImpl) PlaceBet(matchId int32, generatedOdds domain.MatchProbability, currentPot float64) results.Option[domain.BetOrder] {
	return p.betPlacer.Place1x2Bets(matchId, generatedOdds, currentPot)
}
