package backtest

import (
	"errors"
	"fmt"
	"sports-book.com/model"
	"sports-book.com/predict"
	"sports-book.com/predict/goals_predictor"
	"sports-book.com/predict/probability_generator"
	"sports-book.com/util"
	"strconv"
)

func RunBacktests(startYear, endYear int32, pipeline predict.Pipeline) {
	var probabilitiesForCalibration = make(map[model.Match]probability_generator.MatchProbability)
	for i := startYear; i <= endYear; i++ {
		yearProbabilities, err := testPredictSeason(pipeline, i, false)
		if err != nil {
			fmt.Println(err)
			return
		}
		for k, v := range yearProbabilities {
			probabilitiesForCalibration[k] = v
		}
	}

	startSlice := strconv.Itoa(int(startYear))
	endSlice := strconv.Itoa(int(endYear))
	yearTag := fmt.Sprintf("%s_%s", startSlice[len(startSlice)-2:], endSlice[len(endSlice)-2:])

	getCalibration(probabilitiesForCalibration, yearTag)
	plotDistribution(probabilitiesForCalibration, yearTag)
}

func testPredictSeason(pipeline predict.Pipeline, season int32, placeBets bool) (map[model.Match]probability_generator.MatchProbability, error) {

	var winningBets = make(map[model.Match]predict.BetOrder)
	var losingBets = make(map[model.Match]predict.BetOrder)
	var probabilitiesForCalibration = make(map[model.Match]probability_generator.MatchProbability)

	var bank float64 = 100

	matches := util.GetFixtures(season)

	for _, match := range matches {
		customProbabilities, oddsDelta, err := pipeline.PredictMatch(match.HomeTeam, match.AwayTeam, season)
		if errors.Is(err, goals_predictor.ErrNoPreviousData) {
			continue
		}
		if err != nil {
			return nil, err
		}

		probabilitiesForCalibration[match] = customProbabilities

		if placeBets {
			betPlaced := predict.HandleOddsDelta(oddsDelta, match.ID)
			if betPlaced == nil {
				continue
			}
			bank -= betPlaced.Amount
			potentitalReturn := betPlaced.Amount * betPlaced.OddsTaken

			switch betPlaced.Backing {
			case predict.BackHomeWin:
				if match.HomeGoals > match.AwayGoals {
					bank += potentitalReturn
					winningBets[match] = *betPlaced
				} else {
					losingBets[match] = *betPlaced
				}
			case predict.BackDraw:
				if match.HomeGoals == match.AwayGoals {
					bank += potentitalReturn
					winningBets[match] = *betPlaced
				} else {
					losingBets[match] = *betPlaced
				}
			case predict.BackAwayWin:
				if match.AwayGoals > match.HomeGoals {
					bank += potentitalReturn
					winningBets[match] = *betPlaced
				} else {
					losingBets[match] = *betPlaced
				}
			default:
				return nil, errors.New("invalid bet type")
			}
		}

	}

	return probabilitiesForCalibration, nil
}
