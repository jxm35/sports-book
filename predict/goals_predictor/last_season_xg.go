package goals_predictor

import (
	"fmt"
	"sports-book.com/util"
)

type LastSeasonXgGoalPredictor struct{}

func (*LastSeasonXgGoalPredictor) PredictScore(homeTeam, awayTeam, season int32) (float64, float64, error) {

	// calculate standard for the year before
	seasonStats, err := util.GetSeasonDetails(season - 1)
	if err != nil {
		return -1, -1, err
	}
	avgHomeXg := seasonStats.TotalHomexG / 380
	avgAwayXg := seasonStats.TotalAwayxG / 380
	avgHomeGoalsConceded := avgAwayXg
	avgAwayGoalsConceded := avgHomeXg

	// calculate home team's strengths
	homeSeason, err := util.GetTeamSeasonDetails(season-1, homeTeam)
	if err != nil {
		return -1, -1, err
	}
	if homeSeason.XGScoredAtHome == 0 && homeSeason.XGConcededAtHome == 0 {
		return -1, -1, ErrNoPreviousData
	}
	homeAttackStrength := (float64(homeSeason.GoalsScoredAtHome) / float64(19)) / avgHomeXg
	homeDefenseStrength := (float64(homeSeason.GoalsConcededAtHome) / float64(19)) / avgHomeGoalsConceded

	// calculate away team's strengths
	awaySeason, err := util.GetTeamSeasonDetails(season-1, awayTeam)
	if err != nil {
		return -1, -1, err
	}
	if awaySeason.XGScoredAtHome == 0 && awaySeason.XGConcededAtHome == 0 {
		return -1, -1, ErrNoPreviousData
	}
	awayDefenseStrength := (float64(awaySeason.GoalsConcededAway) / float64(19)) / avgAwayGoalsConceded
	awayAttackStrength := (float64(awaySeason.GoalsScoredAway) / float64(19)) / avgAwayXg

	// use strengths to project home and away goals
	projectedHomeGoals := homeAttackStrength * awayDefenseStrength * avgHomeXg
	projectedAwayGoals := awayAttackStrength * homeDefenseStrength * avgAwayXg

	fmt.Printf("%s: %f | %s: %f", homeTeam, projectedHomeGoals, awayTeam, projectedAwayGoals)
	return projectedHomeGoals, projectedAwayGoals, nil
}
