package goals_predictor

import (
	"errors"
	"fmt"

	"sports-book.com/util"
)

var ErrNoPreviousData = errors.New("no previous data for one or both teams")

type LastSeasonResultGoalPredictor struct{}

func (*LastSeasonResultGoalPredictor) PredictScore(homeTeam, awayTeam, season int32, league string) (float64, float64, error) {
	// calculate standard for the year before
	seasonStats, err := util.GetSeasonDetails(season-1, league)
	if err != nil {
		return -1, -1, err
	}
	avgHomeGoals := float64(seasonStats.TotalHG / 380)
	avgAwayGoals := float64(seasonStats.TotalAG / 380)
	avgHomeGoalsConceded := avgAwayGoals
	avgAwayGoalsConceded := avgHomeGoals

	// calculate home team's strengths
	homeSeason, err := util.GetTeamSeasonDetails(season-1, homeTeam)
	if err != nil {
		return -1, -1, err
	}
	if homeSeason.GoalsScoredAtHome == 0 && homeSeason.GoalsConcededAtHome == 0 {
		return -1, -1, ErrNoPreviousData
	}
	homeAttackStrength := (float64(homeSeason.GoalsScoredAtHome) / float64(19)) / avgHomeGoals
	homeDefenseStrength := (float64(homeSeason.GoalsConcededAtHome) / float64(19)) / avgHomeGoalsConceded

	// calculate away team's strengths
	awaySeason, err := util.GetTeamSeasonDetails(season-1, awayTeam)
	if err != nil {
		return -1, -1, err
	}
	if awaySeason.GoalsScoredAtHome == 0 && awaySeason.GoalsConcededAtHome == 0 {
		return -1, -1, ErrNoPreviousData
	}
	awayDefenseStrength := (float64(awaySeason.GoalsConcededAway) / float64(19)) / avgAwayGoalsConceded
	awayAttackStrength := (float64(awaySeason.GoalsScoredAway) / float64(19)) / avgAwayGoals

	// use strengths to project home and away goals
	projectedHomeGoals := homeAttackStrength * awayDefenseStrength * avgHomeGoals
	projectedAwayGoals := awayAttackStrength * homeDefenseStrength * avgAwayGoals

	fmt.Printf("%d: %f | %d: %f", homeTeam, projectedHomeGoals, awayTeam, projectedAwayGoals)
	return projectedHomeGoals, projectedAwayGoals, nil
}
