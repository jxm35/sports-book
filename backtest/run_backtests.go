package backtest

import (
	"errors"
	"fmt"
	"sports-book.com/model"
	"sports-book.com/predict"
	"sports-book.com/predict/domain"
	"sports-book.com/predict/goals_predictor"
	"sports-book.com/util"
	"strconv"
)

func RunBacktests(startYear, endYear int32, league string, pipeline predict.Pipeline, placeBets bool) {
	var probabilitiesForCalibration = make(map[model.Match]domain.MatchProbability)
	bank := float64(100)
	for i := startYear; i <= endYear; i++ {
		yearProbabilities, resultingBank, err := testPredictSeason(pipeline, i, league, placeBets, bank)
		if err != nil {
			fmt.Println(err)
			return
		}
		for k, v := range yearProbabilities {
			probabilitiesForCalibration[k] = v
		}
		bank = resultingBank
	}

	startSlice := strconv.Itoa(int(startYear))
	endSlice := strconv.Itoa(int(endYear))
	yearTag := fmt.Sprintf("%s_%s", startSlice[len(startSlice)-2:], endSlice[len(endSlice)-2:])

	getCalibration(probabilitiesForCalibration, yearTag)
	plotDistribution(probabilitiesForCalibration, yearTag)
}

func testPredictSeason(pipeline predict.Pipeline, season int32, league string, placeBets bool, bank float64) (map[model.Match]domain.MatchProbability, float64, error) {

	var winningBets = make(map[model.Match]domain.BetOrder)
	var losingBets = make(map[model.Match]domain.BetOrder)
	var probabilitiesForCalibration = make(map[model.Match]domain.MatchProbability)

	matches := util.GetFixtures(season, league)
	if len(matches) <= 0 {
		panic("invalid season provided")
	}

	for _, match := range matches {
		customProbabilities, _, err := pipeline.PredictMatch(match.HomeTeam, match.AwayTeam, season, league)
		if errors.Is(err, goals_predictor.ErrNoPreviousData) {
			continue
		}
		if err != nil {
			return nil, -1, err
		}

		probabilitiesForCalibration[match] = customProbabilities

		if placeBets {
			betOp := pipeline.PlaceBet(match.ID, customProbabilities, bank)
			//betPlaced := predict.HandleOddsDelta(oddsDelta, match.ID)
			if betOp.IsNone() {
				continue
			}
			bet := betOp.Value()
			bank -= bet.Amount
			potentitalReturn := bet.Amount * bet.OddsTaken

			switch bet.Backing {
			case domain.BackHomeWin:
				if match.HomeGoals > match.AwayGoals {
					bank += potentitalReturn
					winningBets[match] = bet
				} else {
					losingBets[match] = bet
				}
			case domain.BackDraw:
				if match.HomeGoals == match.AwayGoals {
					bank += potentitalReturn
					winningBets[match] = bet
				} else {
					losingBets[match] = bet
				}
			case domain.BackAwayWin:
				if match.AwayGoals > match.HomeGoals {
					bank += potentitalReturn
					winningBets[match] = bet
				} else {
					losingBets[match] = bet
				}
			default:
				return nil, -1, errors.New("invalid bet type")
			}
		}

	}

	return probabilitiesForCalibration, bank, nil
}
