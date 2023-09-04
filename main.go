package main

import (
	"fmt"

	"sports-book.com/backtest"
	"sports-book.com/predict"
	"sports-book.com/predict/bet_placer"
	"sports-book.com/predict/domain"
	"sports-book.com/predict/goals_predictor"
	"sports-book.com/predict/probability_generator"
	"sports-book.com/util"
)

// 41886 - 22-23
// 37036 - 21-22
// 29415 - 20-21

func main() {
	util.ConnectDB()
	pipeline, err := predict.NewPipelineBuilder().
		SetPredictor(goals_predictor.NewEloGoalsPredictor(0)).
		// SetPredictor(&goals_predictor.LastSeasonXgGoalPredictor{LastXGames: 0}).
		SetProbabilityGenerator(&probability_generator.WeibullOddsGenerator{}).
		SetBetPlacer(&bet_placer.KellyCriterionBetPlacer{
			MaxPercentBetted: 0.2,
			MinOddsDelta:     0.1,
			MaxOddsDelta:     0.3,
		}).
		Build()
	if err != nil {
		panic(err)
	}
	backtest.RunBacktests(2014, 2022, domain.LeagueSerieA, pipeline, true)
}

func testDb() {
	topScorer, err := util.GetTopScorerInSeason(2022)
	fmt.Println(topScorer)
	fmt.Println(err)

	seasonStats, err := util.GetSeasonDetails(2022, "epl")
	fmt.Println(seasonStats)
	fmt.Println(err)

	chelsea := util.GetTeam("Chelsea")
	chelseaSeasonStats, err := util.GetTeamSeasonDetails(2022, chelsea.ID)
	fmt.Println(chelseaSeasonStats)
	fmt.Println(err)
}
