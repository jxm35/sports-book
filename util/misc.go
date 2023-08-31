package util

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"sports-book.com/model"
	"sports-book.com/query"
)

var dbConn *gorm.DB

func xy(db *gorm.DB, err error) int {
	return 5
}

func ConnectDB() error {
	gormDb, err := gorm.Open(mysql.Open("root:password@tcp(127.0.0.1:3306)/sports-book?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		fmt.Println(err)
		return err
	}
	query.SetDefault(gormDb)
	dbConn = gormDb
	return nil
}

func getConn() *gorm.DB {
	if dbConn == nil {
		ConnectDB()
	}
	return dbConn
}

func GetFixtures(season int32, league string) []model.Match {
	var fixtureList []model.Match
	m := query.Match
	c := query.Competition
	m.WithContext(context.Background()).
		Select(m.ALL).
		Join(c, m.Competition.EqCol(c.ID)).
		Where(c.Year.Eq(season)).
		Where(c.Code.Eq(league)).
		Order(m.Date.Desc()).
		Scan(&fixtureList)
	return fixtureList
}

type OddsChecka struct {
	HomeBookie string
	HomeWin    float64
	DrawBookie string
	Draw       float64
	AwayBookie string
	AwayWin    float64
}

func GetBestOddsForMatch(matchId int32) OddsChecka {
	var resp OddsChecka
	o := query.Odds1x2

	o.WithContext(context.Background()).
		Select(o.Bookmaker.As("HomeBookie"), o.HomeWin).
		Where(o.Match.Eq(matchId)).
		Order(o.HomeWin.Desc()).
		Limit(1).
		Scan(&resp)

	o.WithContext(context.Background()).
		Select(o.Bookmaker.As("DrawBookie"), o.Draw).
		Where(o.Match.Eq(matchId)).
		Order(o.Draw.Desc()).
		Limit(1).
		Scan(&resp)

	o.WithContext(context.Background()).
		Select(o.Bookmaker.As("AwayBookie"), o.AwayWin).
		Where(o.Match.Eq(matchId)).
		Order(o.AwayWin.Desc()).
		Limit(1).
		Scan(&resp)
	return resp
}

func GetBestOdds(homeTeam, awayTeam, year int32) OddsChecka {
	var resp OddsChecka
	m := query.Match
	o := query.Odds1x2
	c := query.Competition

	o.WithContext(context.Background()).
		Select(o.Bookmaker.As("HomeBookie"), o.HomeWin).
		Join(m, m.ID.EqCol(o.Match)).
		Join(c, c.ID.EqCol(m.Competition)).
		Where(m.HomeTeam.Eq(homeTeam), m.AwayTeam.Eq(awayTeam), c.Year.Eq(year)).
		Order(o.HomeWin.Desc()).
		Limit(1).
		Scan(&resp)

	o.WithContext(context.Background()).
		Select(o.Bookmaker.As("DrawBookie"), o.Draw).
		Join(m, m.ID.EqCol(o.Match)).
		Join(c, c.ID.EqCol(m.Competition)).
		Where(m.HomeTeam.Eq(homeTeam), m.AwayTeam.Eq(awayTeam), c.Year.Eq(year)).
		Order(o.Draw.Desc()).
		Limit(1).
		Scan(&resp)

	o.WithContext(context.Background()).
		Select(o.Bookmaker.As("AwayBookie"), o.AwayWin).
		Join(m, m.ID.EqCol(o.Match)).
		Join(c, c.ID.EqCol(m.Competition)).
		Where(m.HomeTeam.Eq(homeTeam), m.AwayTeam.Eq(awayTeam), c.Year.Eq(year)).
		Order(o.AwayWin.Desc()).
		Limit(1).
		Scan(&resp)

	return resp
}

func GetTeam(teamName string) model.Team {
	t := query.Team
	var team model.Team
	t.WithContext(context.Background()).
		Select(t.ALL).
		Where(t.Name.Eq(teamName)).Scan(&team)
	return team
}

type SeasonDetails struct {
	TotalHG     int32
	TotalAG     int32
	TotalHomexG float64
	TotalAwayxG float64
}

func GetSeasonDetails(season int32, league string) (SeasonDetails, error) {
	ctx := context.Background()
	m := query.Match
	c := query.Competition
	var res SeasonDetails
	err := m.WithContext(ctx).
		Select(
			m.HomeGoals.Sum().As("TotalHG"),
			m.AwayGoals.Sum().As("TotalAG"),
			m.HomeExpectedGoals.Sum().As("TotalHomexG"),
			m.AwayExpectedGoals.Sum().As("TotalAwayxG")).
		LeftJoin(c, c.ID.EqCol(m.Competition)).
		Where(c.Year.Eq(season)).
		Where(c.Code.Eq(league)).
		Scan(&res)
	return res, err
}

type TeamSeasonDetails struct {
	GoalsScoredAtHome   int32
	GoalsConcededAtHome int32
	XGScoredAtHome      float64
	XGConcededAtHome    float64

	GoalsScoredAway   int32
	GoalsConcededAway int32
	XGScoredAway      float64
	XGConcededAway    float64
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

func GetTeamHomeMatchesBefore(team int32, since time.Time) ([]model.Match, error) {
	ctx := context.Background()
	var res []model.Match
	m := query.Match
	err := m.WithContext(ctx).
		Select(m.ALL).
		Where(
			m.Date.Date().Lt(since),
			m.HomeTeam.Eq(team),
		).
		Order(m.Date.Desc()).
		Scan(&res)
	return res, err
}

func GetTeamAwayMatchesBefore(team int32, since time.Time) ([]model.Match, error) {
	ctx := context.Background()
	var res []model.Match
	m := query.Match
	err := m.WithContext(ctx).
		Select(m.ALL).
		Where(
			m.Date.Date().Lt(since),
			m.AwayTeam.Eq(team),
		).
		Order(m.Date.Desc()).
		Scan(&res)
	return res, err
}

func GetTeamHomeMatchesSince(team int32, since time.Time) ([]model.Match, error) {
	ctx := context.Background()
	var res []model.Match
	m := query.Match
	err := m.WithContext(ctx).
		Select(m.ALL).
		Where(
			m.Date.Date().Gt(since),
			m.HomeTeam.Eq(team),
		).
		Order(m.Date.Desc()).
		Scan(&res)
	return res, err
}

func GetTeamAwayMatchesSince(team int32, since time.Time) ([]model.Match, error) {
	ctx := context.Background()
	var res []model.Match
	m := query.Match
	err := m.WithContext(ctx).
		Select(m.ALL).
		Where(
			m.Date.Date().Gt(since),
			m.AwayTeam.Eq(team),
		).
		Order(m.Date.Desc()).
		Scan(&res)
	return res, err
}

func GetMatchesInSeason(season int32) ([]model.Match, error) {
	m := query.Match
	c := query.Competition
	ctx := context.Background()
	var res []model.Match
	err := m.WithContext(ctx).
		Select(m.ALL).
		Join(c, m.Competition.EqCol(c.ID)).
		Where(c.Year.Eq(season)).
		Scan(&res)
	return res, err
}

func GetTeamSeasonDetails(season, team int32) (TeamSeasonDetails, error) {
	ctx := context.Background()
	m := query.Match
	c := query.Competition

	var res TeamSeasonDetails
	err := m.WithContext(ctx).
		Select(
			m.HomeGoals.Sum().As("GoalsScoredAtHome"),
			m.AwayGoals.Sum().As("GoalsConcededAtHome"),
			m.HomeExpectedGoals.Sum().As("XGScoredAtHome"),
			m.AwayExpectedGoals.Sum().As("XGConcededAtHome")).
		LeftJoin(c, c.ID.EqCol(m.Competition)).
		Where(c.Year.Eq(season), m.HomeTeam.Eq(team)).Scan(&res)
	if err != nil {
		return TeamSeasonDetails{}, err
	}
	err = m.WithContext(ctx).
		Select(
			m.AwayGoals.Sum().As("GoalsScoredAway"),
			m.HomeGoals.Sum().As("GoalsConcededAway"),
			m.AwayExpectedGoals.Sum().As("XGScoredAway"),
			m.HomeExpectedGoals.Sum().As("XGConcededAway")).
		LeftJoin(c, c.ID.EqCol(m.Competition)).
		Where(c.Year.Eq(season), m.AwayTeam.Eq(team)).Scan(&res)
	return res, err
}

func GetPlayersForTeam(team string) {
	ctx := context.Background()
	type Result struct {
		PlayerName string
		Minutes    int32
		TeamName   string
	}

	var res Result
	p := query.Player
	a := query.Appearance
	t := query.Team
	err := p.WithContext(ctx).
		Select(p.Name.As("PlayerName"), a.Minutes, t.Name.As("TeamName")).
		LeftJoin(a, a.Player.EqCol(p.ID)).
		LeftJoin(t, t.ID.EqCol(a.Team)).
		Where(t.Name.Eq(team)).
		Scan(&res)
	fmt.Println(err)
}

type LineupRes struct {
	model.Player
	model.Appearance
}

func GetLineup(team string, day time.Time) []model.Player {
	p := query.Player
	a := query.Appearance
	t := query.Team
	m := query.Match

	var lRes []model.Player
	err := m.WithContext(context.Background()).Select(p.ALL).
		Join(a, a.Match.EqCol(m.ID)).
		Join(p, p.ID.EqCol(a.Player)).
		Join(t, t.ID.EqCol(m.HomeTeam)).
		Where(
			m.Date.Date().Between(getStartOfDay(day), getEndOfDay(day)),
			t.Name.Eq(team),
		).
		Order(a.Minutes.Desc()).
		Limit(11).
		Scan(&lRes)
	fmt.Println(err)
	return lRes
}

type topScorer struct {
	Name           string
	TotalGoals     int32
	XG             float64
	AverageMinutes float64
}

func GetTopScorerInSeason(season int32) (topScorer, error) {
	var ts topScorer

	p := query.Player
	a := query.Appearance
	m := query.Match
	c := query.Competition

	err := a.WithContext(context.Background()).
		Select(
			p.Name,
			a.Goals.Sum().As("TotalGoals"),
			a.ExpectedGoals.Sum().As("XG"),
			a.Minutes.Avg().As("AverageMinutes"),
			a.ExpectedGoals.Count()).
		LeftJoin(m, m.ID.EqCol(a.Match)).
		LeftJoin(p, p.ID.EqCol(a.Player)).
		LeftJoin(c, c.ID.EqCol(m.Competition)).
		Group(a.Player).
		Where(c.Year.Eq(season)).
		Order(a.Goals.Sum().Desc()).
		Scan(&ts)

	return ts, err
}

func getStartOfDay(day time.Time) time.Time {
	return time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.Local)
}

func getEndOfDay(day time.Time) time.Time {
	return time.Date(day.Year(), day.Month(), day.Day(), 23, 59, 0, 0, time.Local)
}
