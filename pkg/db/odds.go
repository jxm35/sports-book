package db

import (
	"context"
	"errors"

	model "sports-book.com/pkg/db_model"
	"sports-book.com/pkg/db_query"
	"sports-book.com/pkg/domain"
)

// GetBestOddsForMatch returns the best odds of found for a given match, as found in the database.
// The odds can each be provided by different bookmakers.
func GetBestOddsForMatch(ctx context.Context, matchId int32) (domain.BookmakerMatchOdds, error) {
	var resp domain.BookmakerMatchOdds
	o := db_query.Odds1x2

	err := o.WithContext(ctx).
		Select(o.Bookmaker.As("HomeBookie"), o.HomeWin).
		Where(o.Match.Eq(matchId)).
		Order(o.HomeWin.Desc()).
		Limit(1).
		Scan(&resp)
	if err != nil {
		return domain.BookmakerMatchOdds{}, err
	}

	err = o.WithContext(ctx).
		Select(o.Bookmaker.As("DrawBookie"), o.Draw).
		Where(o.Match.Eq(matchId)).
		Order(o.Draw.Desc()).
		Limit(1).
		Scan(&resp)
	if err != nil {
		return domain.BookmakerMatchOdds{}, err
	}

	err = o.WithContext(ctx).
		Select(o.Bookmaker.As("AwayBookie"), o.AwayWin).
		Where(o.Match.Eq(matchId)).
		Order(o.AwayWin.Desc()).
		Limit(1).
		Scan(&resp)
	if err != nil {
		return domain.BookmakerMatchOdds{}, err
	}
	return resp, err
}

func SaveOdds(ctx context.Context, odds model.Odds1x2) error {
	o := db_query.Odds1x2
	err := o.WithContext(ctx).Create(&odds)
	return err
}

func SaveBetPlaced(ctx context.Context, bet domain.BetOrder) error {
	b := db_query.BetsPlaced

	oddsTaken, err := getOddsFromBookmakerAndPrice(ctx, bet.BookMaker, bet.OddsTaken, bet.Backing)
	if err != nil {
		return err
	}
	err = b.WithContext(ctx).Create(&model.BetsPlaced{
		MatchID: bet.MatchId,
		Odds:    oddsTaken.ID,
		Amount:  bet.Amount,
	})
	return err
}

func getOddsFromBookmakerAndPrice(ctx context.Context, bookmaker string, oddsTaken float64, backing domain.BackingType) (model.Odds1x2, error) {
	o := db_query.Odds1x2
	var err error
	var res model.Odds1x2
	switch backing {
	case domain.BackHomeWin:
		err = o.WithContext(ctx).Select(o.ALL).Where(o.Bookmaker.Eq(bookmaker), o.HomeWin.Eq(oddsTaken)).Scan(&res)
	case domain.BackDraw:
		err = o.WithContext(ctx).Select(o.ALL).Where(o.Bookmaker.Eq(bookmaker), o.Draw.Eq(oddsTaken)).Scan(&res)
	case domain.BackAwayWin:
		err = o.WithContext(ctx).Select(o.ALL).Where(o.Bookmaker.Eq(bookmaker), o.AwayWin.Eq(oddsTaken)).Scan(&res)
	default:
		return model.Odds1x2{}, errors.New("invalid backing type")
	}
	return res, err
}

/*
	Unused
*/

// getBestOdds returns the best odds of found for a given match, as found in the database.
// The odds can each be provided by different bookmakers.
func getBestOdds(homeTeam, awayTeam, year int32) domain.BookmakerMatchOdds {
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
