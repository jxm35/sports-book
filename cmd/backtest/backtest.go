package main

import (
	"sports-book.com/internal/backtest"
	"sports-book.com/pkg/predict"
	"sports-book.com/pkg/predict/bet_placer"
	"sports-book.com/pkg/predict/domain"
	"sports-book.com/pkg/predict/goals_predictor"
	"sports-book.com/pkg/predict/probability_generator"
	"sports-book.com/pkg/util"
)

func main() {
	util.ConnectDB()
	pipeline, err := predict.NewPipelineBuilder().
		SetPredictor(goals_predictor.NewEloGoalsPredictor(5, 11)).
		// SetPredictor(&goals_predictor.LastSeasonXgGoalPredictor{LastXGames: 0}).
		SetProbabilityGenerator(&probability_generator.WeibullOddsGenerator{}).
		// SetBetPlacer(bet_placer.NewKellyCriterionBetPlacer(0.1, 0.3, 0.05, true)).
		SetBetPlacer(bet_placer.NewFixedAmountBetPlacer(0.1, 0.3, 0.2)).
		Build()
	if err != nil {
		panic(err)
	}
	backtest.RunBacktests(2021, 2022, domain.LeagueEPL, pipeline, true)
}
