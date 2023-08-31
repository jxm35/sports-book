package goals_predictor

import (
	"fmt"

	"sports-book.com/util"
)

type LastSeasonXgGoalPredictor struct{}

func (*LastSeasonXgGoalPredictor) PredictScore(homeTeam, awayTeam, season int32, league string) (float64, float64, error) {
	// calculate standard for the year before
	seasonStats, err := util.GetSeasonDetails(season-1, league)
	if err != nil {
		return -1, -1, err
	}
	if seasonStats.TotalHG == 0 && seasonStats.TotalAG == 0 {
		return -1, -1, ErrNoPreviousData
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
	homeAttackStrength := (homeSeason.XGScoredAtHome / float64(19)) / avgHomeXg
	homeDefenseStrength := (homeSeason.XGConcededAtHome / float64(19)) / avgHomeGoalsConceded

	// calculate away team's strengths
	awaySeason, err := util.GetTeamSeasonDetails(season-1, awayTeam)
	if err != nil {
		return -1, -1, err
	}
	if awaySeason.XGScoredAtHome == 0 && awaySeason.XGConcededAtHome == 0 {
		return -1, -1, ErrNoPreviousData
	}
	awayDefenseStrength := (awaySeason.XGConcededAway / float64(19)) / avgAwayGoalsConceded
	awayAttackStrength := (awaySeason.XGScoredAway / float64(19)) / avgAwayXg

	// use strengths to project home and away goals
	projectedHomeGoals := homeAttackStrength * awayDefenseStrength * avgHomeXg
	projectedAwayGoals := awayAttackStrength * homeDefenseStrength * avgAwayXg

	fmt.Printf("%d: %f | %d: %f", homeTeam, projectedHomeGoals, awayTeam, projectedAwayGoals)
	return projectedHomeGoals, projectedAwayGoals, nil
}
