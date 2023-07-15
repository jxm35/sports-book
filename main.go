package main

import (
	"fmt"
	"sports-book.com/backtest"
	"sports-book.com/predict"
	"sports-book.com/predict/goals_predictor"
	"sports-book.com/predict/probability_generator"
	"sports-book.com/util"
)

// 41886 - 22-23
// 37036 - 21-22
// 29415 - 20-21

func main() {
	//}
	//testDb()
	//save()
	util.ConnectDB()
	pipeline := predict.NewPipelineBuilder().
		SetPredictor(&goals_predictor.LastSeasonXgGoalPredictor{}).
		SetProbabilityGenerator(&probability_generator.WeibullOddsGenerator{}).
		Build()
	backtest.RunBacktests(2015, 2022, pipeline)
	//weib := probability_generator.WeibullOddsGenerator{}
	//results := weib.Generate1x2Probabilities(3.12, 2.09)
	//poiss := probability_generator.PoissonOddsGenerator{}
	//pResults := poiss.Generate1x2Probabilities(3.12, 2.09)
	//fmt.Println(results)
	//fmt.Println(pResults)
}

func testDb() {
	topScorer, err := util.GetTopScorerInSeason(2022)
	fmt.Println(topScorer)
	fmt.Println(err)

	seasonStats, err := util.GetSeasonDetails(2022)
	fmt.Println(seasonStats)
	fmt.Println(err)

	chelsea := util.GetTeam("Chelsea")
	chelseaSeasonStats, err := util.GetTeamSeasonDetails(2022, chelsea.ID)
	fmt.Println(chelseaSeasonStats)
	fmt.Println(err)
}
