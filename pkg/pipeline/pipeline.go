package pipeline

import (
	"errors"
	"fmt"
	"time"

	results "github.com/jxm35/go-results"

	"sports-book.com/pkg/bet_selector"
	"sports-book.com/pkg/db"
	"sports-book.com/pkg/domain"
	"sports-book.com/pkg/probability_generator"
	"sports-book.com/pkg/score_predictor"
)

var ErrTeamNotFound = errors.New("team not found")

type pipelineImpl struct {
	predictor            score_predictor.ScorePredictor
	probabilityGenerator probability_generator.ProbabilityGenerator
	betPlacer            bet_selector.BetSelector
}

func (p *pipelineImpl) PredictMatch(homeTeam, awayTeam, seasonYear int32, league domain.League, date time.Time, matchID int32) (domain.MatchProbability, domain.OddsDelta, error) {
	homeGoalsPredicted, awayGoalsPredicted, err := p.predictor.PredictScore(homeTeam, awayTeam, seasonYear, league, date, matchID)
	if err != nil {
		return domain.MatchProbability{}, domain.OddsDelta{}, err
	}
	matchProbabilities := p.probabilityGenerator.Generate1x2Probabilities(homeGoalsPredicted, awayGoalsPredicted, league)
	fmt.Printf("my probabilities: %+v\n", matchProbabilities)

	bestOdds := db.GetBestOdds(homeTeam, awayTeam, seasonYear)
	fmt.Printf("best odds:%+v\n", bestOdds)

	impliedOdds := domain.MatchProbability{
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
