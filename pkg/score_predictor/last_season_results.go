package score_predictor

import (
	"errors"
	"fmt"
	"time"

	"sports-book.com/pkg/domain"
	"sports-book.com/pkg/entity"
)

var (
	ErrNoPreviousData = errors.New("no previous data for one or both teams")
	ErrInvalidSeason  = errors.New("invalid past season data retrieved")
)

type LastSeasonResultScorePredictor struct{}

func (l *LastSeasonResultScorePredictor) PredictScore(homeTeam, awayTeam, season int32, league domain.League, date time.Time, matchID int32) (float64, float64, error) {
	// calculate standard for the year before
	seasonStats, err := entity.GetSeasonDetails(season-1, league)
	if err != nil {
		return -1, -1, err
	}
	if seasonStats.MatchCount == 0 || seasonStats.MatchCount%2 != 0 {
		return -1, -1, ErrInvalidSeason
	}
	avgHomeGoals := float64(seasonStats.TotalHG / int32(seasonStats.MatchCount))
	avgAwayGoals := float64(seasonStats.TotalAG / int32(seasonStats.MatchCount))
	avgHomeGoalsConceded := avgAwayGoals
	avgAwayGoalsConceded := avgHomeGoals

	// calculate home team's strengths
	homeSeason, err := entity.GetTeamSeasonDetails(season-1, homeTeam)
	if err != nil {
		return -1, -1, err
	}
	if homeSeason.HomeCount == 0 || homeSeason.AwayCount == 0 || homeSeason.HomeCount != homeSeason.AwayCount {
		return -1, -1, ErrInvalidSeason
	}
	if homeSeason.GoalsScoredAtHome == 0 && homeSeason.GoalsConcededAtHome == 0 {
		return -1, -1, ErrNoPreviousData
	}
	homeAttackStrength := (float64(homeSeason.GoalsScoredAtHome) / float64(homeSeason.AwayCount)) / avgHomeGoals
	homeDefenseStrength := (float64(homeSeason.GoalsConcededAtHome) / float64(homeSeason.AwayCount)) / avgHomeGoalsConceded

	// calculate away team's strengths
	awaySeason, err := entity.GetTeamSeasonDetails(season-1, awayTeam)
	if err != nil {
		return -1, -1, err
	}
	if awaySeason.HomeCount == 0 || awaySeason.AwayCount == 0 || awaySeason.HomeCount != awaySeason.AwayCount {
		return -1, -1, ErrInvalidSeason
	}
	if awaySeason.GoalsScoredAtHome == 0 && awaySeason.GoalsConcededAtHome == 0 {
		return -1, -1, ErrNoPreviousData
	}
	awayDefenseStrength := (float64(awaySeason.GoalsConcededAway) / float64(awaySeason.AwayCount)) / avgAwayGoalsConceded
	awayAttackStrength := (float64(awaySeason.GoalsScoredAway) / float64(awaySeason.AwayCount)) / avgAwayGoals

	// use strengths to project home and away goals
	projectedHomeGoals := homeAttackStrength * awayDefenseStrength * avgHomeGoals
	projectedAwayGoals := awayAttackStrength * homeDefenseStrength * avgAwayGoals

	fmt.Printf("%d: %f | %d: %f", homeTeam, projectedHomeGoals, awayTeam, projectedAwayGoals)
	return projectedHomeGoals, projectedAwayGoals, nil
}
