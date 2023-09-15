package odds

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/samber/lo"

	"sports-book.com/pkg/db"
	"sports-book.com/pkg/domain"
	"sports-book.com/pkg/gorm/model"
	"sports-book.com/pkg/notify"
)

type apiRetriever struct{}

// oddsApiNames maps from the name in our database to the name as stored on oddsApi
// todo: finish this
var oddsApiNames = map[string]string{
	"Everton":                 "Everton",
	"Bournemouth":             "Bournemouth",
	"Southampton":             "74",
	"Leicester":               "75",
	"Crystal Palace":          "Crystal Palace",
	"Chelsea":                 "Chelsea",
	"West Ham":                "West Ham United",
	"Tottenham":               "Tottenham Hotspur",
	"Arsenal":                 "Arsenal",
	"Newcastle United":        "Newcastle United",
	"Liverpool":               "Liverpool",
	"Manchester City":         "Manchester City",
	"Manchester United":       "Manchester United",
	"Watford":                 "90",
	"Burnley":                 "Burnley",
	"Huddersfield":            "219",
	"Brighton":                "Brighton and Hove Albion",
	"Cardiff":                 "227",
	"Fulham":                  "Fulham",
	"Wolverhampton Wanderers": "Wolverhampton Wanderers",
	"Aston Villa":             "Aston Villa",
	"Norwich":                 "79",
	"Sheffield United":        "Sheffield United",
	"West Bromwich Albion":    "76",
	"Leeds":                   "245",
	"Brentford":               "Brentford",
	"Nottingham Forest":       "Nottingham Forest",
	"Sunderland":              "77",
	"Swansea":                 "84",
	"Stoke":                   "85",
	"Hull":                    "91",
	"Queens Park Rangers":     "202",
	"Middlesbrough":           "93",
	"Luton":                   "Luton",
}

var ErrOddsNotFount = errors.New("could not find odds for match")

func (a *apiRetriever) GetBestOdds(ctx context.Context, matchId int32, league domain.League) (domain.BookmakerMatchOdds, error) {
	if league == "" {
		notify.GetNotifier().NotifyError("league not provided to api odds retriever")
		league = domain.LeagueEPL
	}
	bestOdds := domain.BookmakerMatchOdds{
		HomeWin: -1,
		Draw:    -1,
		AwayWin: -1,
	}

	match, err := db.GetMatch(ctx, matchId)
	if err != nil {
		return domain.BookmakerMatchOdds{}, err
	}

	home, err := db.GetTeam(ctx, match.HomeTeam)
	if err != nil {
		return bestOdds, err
	}
	away, err := db.GetTeam(ctx, match.AwayTeam)
	if err != nil {
		return bestOdds, err
	}

	odds, err := getOdds(league, match.Date, home.Name, away.Name)
	if err != nil {
		return bestOdds, err
	}

	bestOdds = lo.Reduce(odds.Bookmakers, func(agg domain.BookmakerMatchOdds, bookmakerOdds domain.BookmakerOddsResponse, index int) domain.BookmakerMatchOdds {
		var homeWin, draw, awayWin float64
		for _, market := range bookmakerOdds.Markets {
			if market.Key == "h2h" {
				for _, outcome := range market.Outcomes {
					if outcome.Name == odds.HomeTeam && outcome.Price > agg.HomeWin {
						agg.HomeWin = outcome.Price
						agg.HomeBookie = bookmakerOdds.Title
						homeWin = outcome.Price
					}
					if outcome.Name == odds.AwayTeam && outcome.Price > agg.AwayWin {
						agg.AwayWin = outcome.Price
						agg.AwayBookie = bookmakerOdds.Title
						awayWin = outcome.Price
					}
					if outcome.Name == "Draw" && outcome.Price > agg.Draw {
						agg.Draw = outcome.Price
						agg.DrawBookie = bookmakerOdds.Title
						draw = outcome.Price
					}
				}
			}
		}
		go func() {
			bookmaker := bookmakerOdds.Title
			err := db.SaveOdds(ctx, model.Odds1x2{
				Bookmaker: bookmaker,
				Match:     matchId,
				HomeWin:   homeWin,
				Draw:      draw,
				AwayWin:   awayWin,
			})
			if err != nil {
				notify.GetNotifier().NotifyError(
					fmt.Sprintf("could not save odds: %s", err.Error()),
				)
			}
		}()
		return agg
	}, bestOdds)

	return bestOdds, nil
}

func getOdds(league domain.League, gameDate time.Time, homeName, awayName string) (domain.OddsResponse, error) {
	home, ok := oddsApiNames[homeName]
	if !ok {
		notify.GetNotifier().NotifyError(fmt.Sprintf("tean %s not found in map", homeName))
	}
	away, ok := oddsApiNames[awayName]
	if !ok {
		notify.GetNotifier().NotifyError(fmt.Sprintf("tean %s not found in map", awayName))
	}

	start := getStartOfDay(gameDate).Format(time.RFC3339)
	end := getEndOfDay(gameDate).Format(time.RFC3339)
	apiKey := os.Getenv("ODDS_API_KEY")

	url := fmt.Sprintf(
		"https://api.the-odds-api.com/v4/sports/soccer_%s/odds/?apiKey=%s&regions=uk&markets=h2h&commenceTimeFrom=%s&commenceTimeTo=%s",
		string(league), apiKey, start, end,
	)
	resp, err := http.Get(url)
	if err != nil {
		return domain.OddsResponse{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return domain.OddsResponse{}, err
	}
	var oddsArray []domain.OddsResponse
	if err := json.Unmarshal(body, &oddsArray); err != nil {
		return domain.OddsResponse{}, err
	}

	for _, odds := range oddsArray {
		if odds.HomeTeam == home && odds.AwayTeam == away {
			return odds, nil
		}
	}
	notify.GetNotifier().NotifyError(fmt.Sprintf("could not find teams %s & %s", home, away))
	return domain.OddsResponse{}, ErrOddsNotFount
}

func getStartOfDay(day time.Time) time.Time {
	return time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.Local)
}

func getEndOfDay(day time.Time) time.Time {
	return time.Date(day.Year(), day.Month(), day.Day(), 23, 59, 0, 0, time.Local)
}
