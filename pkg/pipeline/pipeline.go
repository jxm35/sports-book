package pipeline

import (
	"context"
	"errors"
	"time"

	results "github.com/jxm35/go-results"

	"sports-book.com/pkg/bet_selector"
	"sports-book.com/pkg/domain"
	"sports-book.com/pkg/logger"
	"sports-book.com/pkg/probability_generator"
	"sports-book.com/pkg/score_predictor"
)

var ErrTeamNotFound = errors.New("team not found")

type pipelineImpl struct {
	predictor            score_predictor.ScorePredictor
	probabilityGenerator probability_generator.ProbabilityGenerator
	betPlacer            bet_selector.BetSelector
}

func (p *pipelineImpl) PredictMatch(ctx context.Context, homeTeam, awayTeam, seasonYear int32, league domain.League, date time.Time, matchID int32) (domain.MatchProbability, error) {
	homeGoalsPredicted, awayGoalsPredicted, err := p.predictor.PredictScore(ctx, homeTeam, awayTeam, seasonYear, league, date, matchID)
	if err != nil {
		return domain.MatchProbability{}, err
	}
	matchProbabilities := p.probabilityGenerator.Generate1x2Probabilities(homeGoalsPredicted, awayGoalsPredicted, league)
	logger.Info("generated probabilities", "probabilities", matchProbabilities)
	// fmt.Printf("my probabilities: %+v\n", matchProbabilities)

	return matchProbabilities, nil
}

func (p *pipelineImpl) PlaceBet(ctx context.Context, matchId int32, generatedOdds domain.MatchProbability, currentPot float64) results.Option[domain.BetOrder] {
	return p.betPlacer.Place1x2Bets(ctx, matchId, generatedOdds, currentPot)
}

func getOddsDelta(impliedOdds domain.MatchProbability, bookmakerOdds domain.BookmakerMatchOdds) domain.OddsDelta {
	delta := domain.OddsDelta{
		HomeBookie:   bookmakerOdds.HomeBookie,
		HomeWinDelta: bookmakerOdds.HomeWin - impliedOdds.HomeWin,
		HomeWinOdds:  bookmakerOdds.HomeWin,

		DrawBookie: bookmakerOdds.DrawBookie,
		DrawDelta:  bookmakerOdds.Draw - impliedOdds.Draw,
		DrawOdds:   bookmakerOdds.Draw,

		AwayBookie:   bookmakerOdds.AwayBookie,
		AwayWinDelta: bookmakerOdds.AwayWin - impliedOdds.AwayWin,
		AwayWinOdds:  bookmakerOdds.AwayWin,
	}
	return delta
}
