package score_predictor

import (
	"context"
	"fmt"
	"time"

	"github.com/samber/lo"

	"sports-book.com/pkg/db"
	"sports-book.com/pkg/domain"
	"sports-book.com/pkg/gorm/model"
)

type LastSeasonXgScorePredictor struct {
	LastXGames int
}

func (l *LastSeasonXgScorePredictor) PredictScore(ctx context.Context, homeTeam, awayTeam, season int32, league domain.League, date time.Time, matchID int32) (float64, float64, error) {
	// calculate standard for the year before
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
	homeAvgxG := homeSeason.XGScoredAtHome / float64(homeSeason.AwayCount)
	homeAvgxGConceded := homeSeason.XGConcededAtHome / float64(homeSeason.AwayCount)
	if l.LastXGames != 0 {
		homeLastXGames, err := db.GetHomeLastXMatches(homeTeam, season, date, l.LastXGames)
		if err != nil {
			return -1, -1, err
		}
		if len(homeLastXGames) > 0 {
			lastGamesHomexG := lo.Reduce(homeLastXGames, func(agg float64, item model.Match, index int) float64 {
				return agg + item.HomeExpectedGoals
			}, 0.0)
			lastGamesAwayxG := lo.Reduce(homeLastXGames, func(agg float64, item model.Match, index int) float64 {
				return agg + item.AwayExpectedGoals
			}, 0.0)
			homeAvgxG = (homeAvgxG + (lastGamesHomexG / float64(len(homeLastXGames)))) / 2
			homeAvgxGConceded = (homeAvgxGConceded + (lastGamesAwayxG / float64(len(homeLastXGames)))) / 2
		}
	}

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

	awayAvgxG := awaySeason.XGScoredAway / float64(awaySeason.AwayCount)
	awayAvgxGConceded := awaySeason.XGConcededAway / float64(awaySeason.AwayCount)
	if l.LastXGames != 0 {
		awayLastXGames, err := db.GetAwayLastXMatches(awayTeam, season, date, l.LastXGames)
		if err != nil {
			return -1, -1, err
		}
		if len(awayLastXGames) > 0 {
			lastGamesxGScored := lo.Reduce(awayLastXGames, func(agg float64, item model.Match, index int) float64 {
				return agg + item.AwayExpectedGoals
			}, 0.0)
			lastGamesxGConceded := lo.Reduce(awayLastXGames, func(agg float64, item model.Match, index int) float64 {
				return agg + item.HomeExpectedGoals
			}, 0.0)
			awayAvgxG = (awayAvgxG + (lastGamesxGScored / float64(len(awayLastXGames)))) / 2
			awayAvgxGConceded = (awayAvgxGConceded + (lastGamesxGConceded / float64(len(awayLastXGames)))) / 2
		}
	}

	awayDefenseStrength := awayAvgxGConceded / avgAwayGoalsConceded
	awayAttackStrength := awayAvgxG / avgAwayXg

	// use strengths to project home and away goals
	projectedHomeGoals := homeAttackStrength * awayDefenseStrength * avgHomeXg
	projectedAwayGoals := awayAttackStrength * homeDefenseStrength * avgAwayXg

	fmt.Printf("%d: %f | %d: %f", homeTeam, projectedHomeGoals, awayTeam, projectedAwayGoals)
	return projectedHomeGoals, projectedAwayGoals, nil
}
