package main

import (
	"fmt"
	"sports-book.com/util"
)

// 41886 - 22-23
// 37036 - 21-22
// 29415 - 20-21

func main() {
	//}
	//testDb()
	//save()
	//util.ConnectDB()
	//pipeline, err := predict.NewPipelineBuilder().
	//	SetPredictor(&goals_predictor.LastSeasonXgGoalPredictor{}).
	//	SetProbabilityGenerator(&probability_generator.WeibullOddsGenerator{}).
	//	SetBetPlacer(&bet_placer.KellyCriterionBetPlacer{
	//		MaxPercentBetted: 0.2,
	//		MinOddsDelta:     0.1,
	//		MaxOddsDelta:     0.3,
	//	}).
	//	Build()
	//if err != nil {
	//	panic(err)
	//}
	//backtest.RunBacktests(2018, 2022, pipeline)

	//probability_generator.FindFranksValue()
	//goals_predictor.FindDecayFactor()
	//probability_generator.FindWeibullShapes()

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
