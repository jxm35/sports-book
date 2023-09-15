package score_predictor

import (
	"context"
	"errors"
	"fmt"
	"time"

	"sports-book.com/pkg/db"
	"sports-book.com/pkg/domain"
	"sports-book.com/pkg/gorm/model"
	"sports-book.com/pkg/notify"
	"sports-book.com/pkg/predictions"
)

type eloGoalsPredictor struct {
	predictions               predictions.PredictionHandler
	FormTimespan              int
	chanceConversionWeighting int
}

func NewEloGoalsPredictor(formTimespan, chanceConversionRating int) ScorePredictor {
	return &eloGoalsPredictor{
		predictions:               predictions.GetPredictionHandler(),
		FormTimespan:              formTimespan,
		chanceConversionWeighting: chanceConversionRating,
	}
}

func (e *eloGoalsPredictor) PredictScore(ctx context.Context, homeTeam, awayTeam, season int32, league domain.League, date time.Time, matchID int32) (float64, float64, error) {
	seasonStats, err := db.GetSeasonDetails(season-1, league)
	if err != nil {
		return -1, -1, err
	}
	if seasonStats.TotalHG == 0 && seasonStats.TotalAG == 0 {
		return -1, -1, ErrNoPreviousData
	}
	if seasonStats.MatchCount == 0 || seasonStats.MatchCount%2 != 0 {
		return -1, -1, ErrInvalidSeason
	}
	defer func() {
		if _, predictionErr := e.predictions.GetPrediction(ctx, matchID); errors.Is(predictionErr, predictions.ErrPredictionNotFound) {
			predictionErr = e.predictions.SavePrediction(ctx, matchID, domain.Prediction{
				HomexG: -1,
				AwayxG: -1,
			})
			if predictionErr != nil {
				notify.GetNotifier().NotifyError(predictionErr.Error())
			}
		}
	}()
	avgHomeXg := seasonStats.TotalHomexG / float64(seasonStats.MatchCount)
	avgAwayXg := seasonStats.TotalAwayxG / float64(seasonStats.MatchCount)
	avgHomeGoalsConceded := avgAwayXg
	avgAwayGoalsConceded := avgHomeXg

	// calculate home team's strengths
	homeSeason, err := db.GetTeamSeasonDetails(season-1, homeTeam)
	if err != nil {
		return -1, -1, err
	}
	if homeSeason.XGScoredAtHome == 0 && homeSeason.XGConcededAtHome == 0 {
		return -1, -1, ErrNoPreviousData
	}
	if homeSeason.HomeCount == 0 || homeSeason.AwayCount == 0 || homeSeason.HomeCount != homeSeason.AwayCount {
		return -1, -1, ErrInvalidSeason
	}

	lastGames, err := db.GetHomeLastXMatches(homeTeam, season, date, e.FormTimespan)
	if err != nil {
		return -1, -1, err
	}
	aForm, dForm, err := e.calcFormHome(ctx, homeTeam, lastGames)
	if err != nil {
		return -1, -1, err
	}
	homeChanceConversion := e.calcGoalConversion(homeSeason, homeSeason.HomeCount+homeSeason.AwayCount)

	homeAvgxG := homeSeason.XGScoredAtHome / float64(homeSeason.AwayCount)
	homeAvgxGConceded := homeSeason.XGConcededAtHome / float64(homeSeason.AwayCount)
	homeAvgxG += aForm
	if e.chanceConversionWeighting != 0 {
		homeAvgxG += homeChanceConversion / float64(e.chanceConversionWeighting)
	}
	homeAvgxGConceded += dForm
	homeAttackStrength := homeAvgxG / avgHomeXg
	homeDefenseStrength := homeAvgxGConceded / avgHomeGoalsConceded

	// calculate away team's strengths
	awaySeason, err := db.GetTeamSeasonDetails(season-1, awayTeam)
	if err != nil {
		return -1, -1, err
	}
	if awaySeason.XGScoredAtHome == 0 && awaySeason.XGConcededAtHome == 0 {
		return -1, -1, ErrNoPreviousData
	}
	if awaySeason.HomeCount == 0 || awaySeason.AwayCount == 0 || awaySeason.HomeCount != awaySeason.AwayCount {
		return -1, -1, ErrInvalidSeason
	}

	lastGames, err = db.GetAwayLastXMatches(awayTeam, season, date, e.FormTimespan)
	if err != nil {
		return -1, -1, err
	}
	aForm, dForm, err = e.calcFormAway(ctx, awayTeam, lastGames)
	if err != nil {
		return -1, -1, err
	}
	awayChanceConversion := e.calcGoalConversion(awaySeason, awaySeason.HomeCount+awaySeason.AwayCount)

	awayAvgxG := awaySeason.XGScoredAway / float64(awaySeason.AwayCount)
	awayAvgxGConceded := awaySeason.XGConcededAway / float64(awaySeason.AwayCount)
	awayAvgxG += aForm
	if e.chanceConversionWeighting != 0 {
		awayAvgxG += awayChanceConversion / float64(e.chanceConversionWeighting)
	}
	awayAvgxGConceded += dForm
	awayDefenseStrength := awayAvgxGConceded / avgAwayGoalsConceded
	awayAttackStrength := awayAvgxG / avgAwayXg

	// use strengths to project home and away goals
	projectedHomeGoals := homeAttackStrength * awayDefenseStrength * avgHomeXg
	projectedAwayGoals := awayAttackStrength * homeDefenseStrength * avgAwayXg

	fmt.Printf("%d: %f | %d: %f", homeTeam, projectedHomeGoals, awayTeam, projectedAwayGoals)
	err = e.predictions.SavePrediction(ctx, matchID, domain.Prediction{
		HomexG: projectedHomeGoals,
		AwayxG: projectedAwayGoals,
	})
	return projectedHomeGoals, projectedAwayGoals, err
}

// additive or multiplicative?
// calculate home and away separately or together?

func (e *eloGoalsPredictor) calcFormHome(ctx context.Context, team int32, matches []model.Match) (float64, float64, error) {
	attackForm := 0.0  // how many more goals they score than expected
	defenseForm := 0.0 // how many more goals they concede than expected
	count := 0
	for i := 0; i < len(matches); i++ {
		pred, err := e.predictions.GetPrediction(ctx, matches[i].ID)
		if err != nil {
			return -1, -1, err
		}
		if pred.AwayxG == -1 && pred.HomexG == -1 {
			continue
		}
		attackForm = attackForm + (matches[i].HomeExpectedGoals - pred.HomexG)
		defenseForm = defenseForm + (matches[i].AwayExpectedGoals - pred.AwayxG)
		count++
	}
	if count < 4 {
		return 0, 0, nil
	}
	attackForm = attackForm / float64(count)
	defenseForm = defenseForm / float64(count)

	return attackForm, defenseForm, nil
}

func (e *eloGoalsPredictor) calcFormAway(ctx context.Context, team int32, matches []model.Match) (float64, float64, error) {
	attackForm := 0.0
	defenseForm := 0.0
	count := 0
	for i := 0; i < len(matches); i++ {
		pred, err := e.predictions.GetPrediction(ctx, matches[i].ID)
		if err != nil {
			return -1, -1, err
		}
		if pred.AwayxG == -1 && pred.HomexG == -1 {
			continue
		}
		attackForm = attackForm + (matches[i].AwayExpectedGoals - pred.AwayxG)
		defenseForm = defenseForm + (matches[i].HomeExpectedGoals - pred.HomexG)
		count++
	}
	if count < 4 {
		return 0, 0, nil
	}
	attackForm = attackForm / float64(count)
	defenseForm = defenseForm / float64(count)
	return attackForm, defenseForm, nil
}

func (e *eloGoalsPredictor) calcForm(ctx context.Context, team int32, matches []model.Match) (float64, float64, error) {
	attackForm := 0.0
	defenseForm := 0.0
	count := 0
	for _, match := range matches {
		pred, err := e.predictions.GetPrediction(ctx, match.ID)
		if err != nil {
			return -1, -1, err
		}
		if pred.AwayxG == -1 && pred.HomexG == -1 {
			continue
		}
		if match.HomeTeam == team {
			attackForm = attackForm + (match.HomeExpectedGoals - pred.HomexG)
			defenseForm = defenseForm + (match.AwayExpectedGoals - pred.AwayxG)
		} else if match.AwayTeam == team {
			attackForm = attackForm + (match.AwayExpectedGoals - pred.AwayxG)
			defenseForm = defenseForm + (match.HomeExpectedGoals - pred.HomexG)
		} else {
			panic("team not in match")
		}
		count++
	}
	if count < 4 {
		return 0, 0, nil
	}
	attackForm = attackForm / float64(count)
	defenseForm = defenseForm / float64(count)
	return attackForm, defenseForm, nil
}

func (e *eloGoalsPredictor) calcGoalConversion(lastSeason domain.TeamSeasonDetails, matchesPlayed int) float64 {
	diff := float64(lastSeason.GoalsScoredAtHome+lastSeason.GoalsScoredAway) - (lastSeason.XGScoredAtHome + lastSeason.XGScoredAway)
	return diff / float64(matchesPlayed)
}
