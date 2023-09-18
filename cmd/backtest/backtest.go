package main

import (
	"sports-book.com/internal/backtest"
	"sports-book.com/pkg/db"
	"sports-book.com/pkg/domain"
	"sports-book.com/pkg/logger"
	"sports-book.com/pkg/pipeline"
)

func main() {
	_, err := db.Connect()
	if err != nil {
		panic(err)
	}
	logger.InitialiseDevLogger()
	// p, err := pipeline.NewPipelineBuilder().
	//	SetPredictor(score_predictor.NewEloGoalsPredictor(5, 11)).
	//	// SetPredictor(&goals_predictor.LastSeasonXgGoalPredictor{LastXGames: 0}).
	//	SetProbabilityGenerator(&probability_generator.WeibullOddsGenerator{}).
	//	SetBetPlacer(bet_selector.NewKellyCriterionBetSelector(0.1, 0.3, 0.05, true)).
	//	// SetBetPlacer(bet_selector.NewFixedAmountBetSelector(0.1, 0.3, 0.2)).
	//	Build()
	p, err := pipeline.NewPipelineFromConfig()
	if err != nil {
		panic(err)
	}
	backtest.RunBacktests(2014, 2023, domain.LeagueEPL, p, true)
}
