package goals_predictor

import (
	"fmt"
	"github.com/samber/lo"
	"sports-book.com/model"
	"time"

	"sports-book.com/util"
)

type LastSeasonXgGoalPredictor struct {
	LastXGames int
}

func (l *LastSeasonXgGoalPredictor) PredictScore(homeTeam, awayTeam, season int32, league string, date time.Time, matchID int32) (float64, float64, error) {
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
	homeAvgxG := homeSeason.XGScoredAtHome / float64(19)
	homeAvgxGConceded := homeSeason.XGConcededAtHome / float64(19)
	if l.LastXGames != 0 {
		homeLastXGames, err := util.GetHomeLastXGames(homeTeam, season, date, l.LastXGames)
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
	awaySeason, err := util.GetTeamSeasonDetails(season-1, awayTeam)
	if err != nil {
		return -1, -1, err
	}
	if awaySeason.XGScoredAtHome == 0 && awaySeason.XGConcededAtHome == 0 {
		return -1, -1, ErrNoPreviousData
	}

	awayAvgxG := awaySeason.XGScoredAway / float64(19)
	awayAvgxGConceded := awaySeason.XGConcededAway / float64(19)
	if l.LastXGames != 0 {
		awayLastXGames, err := util.GetAwayLastXGames(awayTeam, season, date, l.LastXGames)
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
