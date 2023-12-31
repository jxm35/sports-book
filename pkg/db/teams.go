package db

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"sports-book.com/pkg/domain"
	"sports-book.com/pkg/gorm/model"
	"sports-book.com/pkg/gorm/query"
)

// GetTeamSeasonDetails returns details of the supplied team's goals and xG for the season that starts with the
// specified year.
func GetTeamSeasonDetails(seasonYear, teamId int32) (domain.TeamSeasonDetails, error) {
	ctx := context.Background()
	m := query.Match
	c := query.Competition

	var res domain.TeamSeasonDetails
	err := m.WithContext(ctx).
		Select(
			m.ID.Count().As("HomeCount"),
			m.HomeGoals.Sum().As("GoalsScoredAtHome"),
			m.AwayGoals.Sum().As("GoalsConcededAtHome"),
			m.HomeExpectedGoals.Sum().As("XGScoredAtHome"),
			m.AwayExpectedGoals.Sum().As("XGConcededAtHome")).
		LeftJoin(c, c.ID.EqCol(m.Competition)).
		Where(c.Year.Eq(seasonYear), m.HomeTeam.Eq(teamId)).Scan(&res)
	if err != nil {
		return domain.TeamSeasonDetails{}, err
	}
	err = m.WithContext(ctx).
		Select(
			m.ID.Count().As("AwayCount"),
			m.AwayGoals.Sum().As("GoalsScoredAway"),
			m.HomeGoals.Sum().As("GoalsConcededAway"),
			m.AwayExpectedGoals.Sum().As("XGScoredAway"),
			m.HomeExpectedGoals.Sum().As("XGConcededAway")).
		LeftJoin(c, c.ID.EqCol(m.Competition)).
		Where(c.Year.Eq(seasonYear), m.AwayTeam.Eq(teamId)).Scan(&res)
	return res, err
}

// GetHomeLastXMatches returns a slice of (up to) the last numMatches matches that the team has played at home.
// These are ordered from most recent to oldest.
func GetHomeLastXMatches(teamId, seasonYear int32, date time.Time, numGames int) ([]model.Match, error) {
	ctx := context.Background()
	m := query.Match
	c := query.Competition

	var res []model.Match
	err := m.WithContext(ctx).
		Select(m.ALL).
		LeftJoin(c, c.ID.EqCol(m.Competition)).
		Where(
			c.Year.Eq(seasonYear),
			m.HomeTeam.Eq(teamId),
			m.Date.Lt(getStartOfDay(date)),
		).
		Order(m.Date.Desc()).
		Limit(numGames).
		Scan(&res)
	return res, err
}

// GetAwayLastXMatches returns a slice of (up to) the last numMatches matches that the team has played at away from home.
// These are ordered from most recent to oldest.
func GetAwayLastXMatches(team, season int32, date time.Time, numGames int) ([]model.Match, error) {
	ctx := context.Background()
	m := query.Match
	c := query.Competition

	var res []model.Match
	err := m.WithContext(ctx).
		Select(m.ALL).
		LeftJoin(c, c.ID.EqCol(m.Competition)).
		Where(
			c.Year.Eq(season),
			m.AwayTeam.Eq(team),
			m.Date.Lt(getStartOfDay(date)),
		).
		Order(m.Date.Desc()).
		Limit(numGames).
		Scan(&res)
	return res, err
}

// GetTeamHomeMatchesSince returns a slice of the  matches that the team has played at home since the given date.
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

// GetTeamAwayMatchesSince returns a slice of the  matches that the team has played at away since the given date.
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

// GetTeam retrieves a team from the database given it's id
func GetTeam(ctx context.Context, teamId int32) (model.Team, error) {
	t := query.Team
	team, err := t.WithContext(ctx).
		Select(t.ALL).
		Where(t.ID.Eq(teamId)).
		First()
	if err != nil {
		return model.Team{}, err
	}
	return *team, err
}

// GetTeamByName retrieves a team from the database given it's name
func GetTeamByName(name string) (model.Team, error) {
	t := query.Team
	team, err := t.WithContext(context.Background()).
		Select(t.ALL).
		Where(t.Name.Eq(name)).
		First()
	if err != nil {
		return model.Team{}, err
	}
	return *team, err
}

// GetTeamByUsId retrieves a team from the database given it's Understat id
func GetTeamByUsId(usId string) (model.Team, error) {
	id, err := strconv.Atoi(usId)
	if err != nil {
		return model.Team{}, err
	}
	t := query.Team
	team, err := t.WithContext(context.Background()).
		Select(t.ALL).
		Where(t.UsID.Eq(int32(id))).
		First()
	if err != nil {
		return model.Team{}, err
	}
	return *team, err
}

/*
	Unused
*/

func ListTeamHomeMatchesBefore(team int32, since time.Time) ([]model.Match, error) {
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

func ListTeamAwayMatchesBefore(team int32, since time.Time) ([]model.Match, error) {
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
