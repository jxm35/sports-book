package db

import (
	"context"

	"sports-book.com/pkg/db_query"
	"sports-book.com/pkg/domain"
)

// GetBestOddsForMatch returns the best odds of found for a given match, as found in the database.
// The odds can each be provided by different bookmakers.
func GetBestOddsForMatch(matchId int32) domain.BookmakerMatchOdds {
	var resp domain.BookmakerMatchOdds
	o := db_query.Odds1x2

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

// GetBestOdds returns the best odds of found for a given match, as found in the database.
// The odds can each be provided by different bookmakers.
func GetBestOdds(homeTeam, awayTeam, year int32) domain.BookmakerMatchOdds {
	var resp domain.BookmakerMatchOdds
	m := db_query.Match
	o := db_query.Odds1x2
	c := db_query.Competition

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
