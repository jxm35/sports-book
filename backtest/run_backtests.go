package backtest

import (
	"errors"
	"fmt"
	"strconv"

	"sports-book.com/model"
	"sports-book.com/predict"
	"sports-book.com/predict/domain"
	"sports-book.com/predict/goals_predictor"
	"sports-book.com/util"
)

func RunBacktests(startYear, endYear int32, league string, pipeline predict.Pipeline, placeBets bool) {
	probabilitiesForCalibration := make(map[model.Match]domain.MatchProbability)
	betResults := make([]betResult, 0)
	bank := float64(100)
	for i := startYear; i <= endYear; i++ {
		yearProbabilities, resultingBank, bets, err := testPredictSeason(pipeline, i, league, placeBets, bank)
		if err != nil {
			fmt.Println(err)
			return
		}
		betResults = append(betResults, bets...)
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
	if placeBets {
		plotBetDistribution(betResults, yearTag)
		plotBetWinningsDistribution(betResults, yearTag)
		plotBetsPlacedDistribution(betResults, yearTag)
	}
}

type betResult struct {
	domain.BetOrder
	domain.MatchProbability
	AmountWon float64
	Won       bool
}

func testPredictSeason(pipeline predict.Pipeline, season int32, league string, placeBets bool, bank float64) (map[model.Match]domain.MatchProbability, float64, []betResult, error) {
	winningBets := make(map[model.Match]domain.BetOrder)
	losingBets := make(map[model.Match]domain.BetOrder)
	probabilitiesForCalibration := make(map[model.Match]domain.MatchProbability)
	betsPlaced := make([]betResult, 0)

	matches := util.GetFixtures(season, league)
	if len(matches) <= 0 {
		panic("invalid season provided")
	}

	for _, match := range matches {
		customProbabilities, _, err := pipeline.PredictMatch(match.HomeTeam, match.AwayTeam, season, league, match.Date, match.ID)
		if errors.Is(err, goals_predictor.ErrNoPreviousData) {
			continue
		}
		if err != nil {
			return nil, -1, nil, err
		}

		probabilitiesForCalibration[match] = customProbabilities

		if placeBets {
			betOp := pipeline.PlaceBet(match.ID, customProbabilities, bank)
			// betPlaced := predict.HandleOddsDelta(oddsDelta, match.ID)
			if betOp.IsNone() {
				continue
			}
			bet := betOp.Value()
			bank -= bet.Amount
			potentitalReturn := bet.Amount * bet.OddsTaken
			betWon := false

			switch bet.Backing {
			case domain.BackHomeWin:
				if match.HomeGoals > match.AwayGoals {
					bank += potentitalReturn
					winningBets[match] = bet
					betWon = true
				} else {
					losingBets[match] = bet
				}
			case domain.BackDraw:
				if match.HomeGoals == match.AwayGoals {
					bank += potentitalReturn
					winningBets[match] = bet
					betWon = true
				} else {
					losingBets[match] = bet
				}
			case domain.BackAwayWin:
				if match.AwayGoals > match.HomeGoals {
					bank += potentitalReturn
					winningBets[match] = bet
					betWon = true
				} else {
					losingBets[match] = bet
				}
			default:
				return nil, -1, nil, errors.New("invalid bet type")
			}
			betsPlaced = append(betsPlaced, betResult{
				BetOrder:         bet,
				MatchProbability: customProbabilities,
				Won:              betWon,
			})
		}

	}

	return probabilitiesForCalibration, bank, betsPlaced, nil
}
