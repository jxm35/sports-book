package entity

import (
	"context"
	"time"

	model "sports-book.com/pkg/db_model"
	"sports-book.com/pkg/db_query"
	"sports-book.com/pkg/domain"
)

// ListFixtures returns all the fixtures stored in the database for the given league for the league starting in the given year.
func ListFixtures(seasonYear int32, league domain.League) ([]model.Match, error) {
	var fixtureList []model.Match
	m := db_query.Match
	c := db_query.Competition
	err := m.WithContext(context.Background()).
		Select(m.ALL).
		Join(c, m.Competition.EqCol(c.ID)).
		Where(c.Year.Eq(seasonYear)).
		Where(c.Code.Eq(string(league))).
		Order(m.Date).
		Scan(&fixtureList)
	return fixtureList, err
}

// GetSeasonDetails returns information about goals and xG scored home and away throughout the season.
func GetSeasonDetails(seasonYear int32, league domain.League) (domain.SeasonDetails, error) {
	ctx := context.Background()
	m := db_query.Match
	c := db_query.Competition
	var res domain.SeasonDetails
	err := m.WithContext(ctx).
		Select(
			m.ID.Count().As("MatchCount"),
			m.HomeGoals.Sum().As("TotalHG"),
			m.AwayGoals.Sum().As("TotalAG"),
			m.HomeExpectedGoals.Sum().As("TotalHomexG"),
			m.AwayExpectedGoals.Sum().As("TotalAwayxG")).
		LeftJoin(c, c.ID.EqCol(m.Competition)).
		Where(c.Year.Eq(seasonYear)).
		Where(c.Code.Eq(string(league))).
		Scan(&res)
	return res, err
}

func GetMatch(matchId int32) (model.Match, error) {
	ctx := context.Background()
	m := db_query.Match
	var res model.Match
	err := m.WithContext(ctx).
		Select(m.ALL).
		Where(m.ID.Eq(matchId)).
		Scan(&res)
	return res, err
}

/*
	unused functions
*/

func GetLastXGames(team, season int32, date time.Time, numGames int) ([]model.Match, error) {
	ctx := context.Background()
	m := db_query.Match
	c := db_query.Competition

	var res []model.Match
	err := m.WithContext(ctx).
		Select(m.ALL).
		LeftJoin(c, c.ID.EqCol(m.Competition)).
		Where(
			c.Year.Eq(season),
			m.HomeTeam.Eq(team),
			m.Date.Lt(getStartOfDay(date)),
		).
		Or(
			c.Year.Eq(season),
			m.AwayTeam.Eq(team),
			m.Date.Lt(getStartOfDay(date)),
		).
		Order(m.Date.Desc()).
		Limit(numGames).
		Scan(&res)
	return res, err
}

func GetHomexGVariance(season int32) float64 {
	db := getConn()
	rawSql := "SELECT VARIANCE(home_expected_goals) FROM `match` m WHERE m.competition = ?;"
	var variance float64
	db.Raw(rawSql, season).Row().Scan(&variance)
	return variance
}

func GetAwayxGVariance(season int32) float64 {
	db := getConn()
	rawSql := "SELECT VARIANCE(away_expected_goals) FROM `match` m WHERE m.competition = ?;"
	var variance float64
	db.Raw(rawSql, season).Row().Scan(&variance)
	return variance
}
