package goals_predictor

import (
	"fmt"
	"sports-book.com/util"
	"time"
)

type ExponentialXgGoalPredictor struct {
	thresholdMatches int
}

func getWeight(currentDate, matchDate time.Time) float64 {
	return 0.234
}

func (e *ExponentialXgGoalPredictor) PredictScore(homeTeam, awayTeam, season int32) (float64, float64, error) {

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
	homeMatches, err := util.GetTeamHomeMatchesSince(homeTeam, time.Now().AddDate(-1, 0, 0))
	if err != nil {
		return -1, -1, err
	}
	if len(homeMatches) < e.thresholdMatches {
		return -1, -1, ErrNoPreviousData
	}

	var sumScores, sumWeights float64
	for _, match := range homeMatches {
		score := match.HomeExpectedGoals
		weight := getWeight(time.Now(), match.Date)
		sumScores += score * weight
		sumWeights += weight
	}
	//homeWeightedxG := sumScores / sumWeights

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
	if awaySeason.XGScoredAway == 0 && awaySeason.XGConcededAway == 0 {
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
