package main

import (
	"sports-book.com/internal/backtest"
	"sports-book.com/pkg/bet_selector"
	"sports-book.com/pkg/domain"
	"sports-book.com/pkg/entity"
	"sports-book.com/pkg/pipeline"
	"sports-book.com/pkg/probability_generator"
	"sports-book.com/pkg/score_predictor"
)

func main() {
	entity.ConnectDB()
	pipeline, err := pipeline.NewPipelineBuilder().
		SetPredictor(score_predictor.NewEloGoalsPredictor(5, 11)).
		// SetPredictor(&goals_predictor.LastSeasonXgGoalPredictor{LastXGames: 0}).
		SetProbabilityGenerator(&probability_generator.WeibullOddsGenerator{}).
		// SetBetPlacer(bet_placer.NewKellyCriterionBetPlacer(0.1, 0.3, 0.05, true)).
		SetBetPlacer(bet_selector.NewFixedAmountBetSelector(0.1, 0.3, 0.2)).
		Build()
	if err != nil {
		panic(err)
	}
	backtest.RunBacktests(2021, 2022, domain.LeagueEPL, pipeline, true)
}
