package domain

type (
	BackingType string
	League      string
)

const (
	BackHomeWin BackingType = "BackHomeWin"
	BackDraw    BackingType = "BackDraw"
	BackAwayWin BackingType = "BackAwayWin"

	LeagueEPL        League = "epl"
	LeagueLaLiga     League = "la_liga"
	LeagueBundesliga League = "bundesliga"
	LeagueSerieA     League = "serie_a"
)

type SeasonDetails struct {
	TotalHG     int32
	TotalAG     int32
	TotalHomexG float64
	TotalAwayxG float64
	MatchCount  int
}

type TeamSeasonDetails struct {
	HomeCount           int
	GoalsScoredAtHome   int32
	GoalsConcededAtHome int32
	XGScoredAtHome      float64
	XGConcededAtHome    float64

	AwayCount         int
	GoalsScoredAway   int32
	GoalsConcededAway int32
	XGScoredAway      float64
	XGConcededAway    float64
}

type BookmakerMatchOdds struct {
	HomeBookie string
	HomeWin    float64

	DrawBookie string
	Draw       float64

	AwayBookie string
	AwayWin    float64
}

type MatchProbability struct {
	HomeWin float64
	Draw    float64
	AwayWin float64
}

type BetOrder struct {
	MatchId              int32
	Backing              BackingType
	BookMaker            string
	OddsTaken            float64
	Amount               float64
	PredictedProbability float64
}

type OddsDelta struct {
	HomeBookie   string
	HomeWinDelta float64
	HomeWinOdds  float64

	DrawBookie string
	DrawDelta  float64
	DrawOdds   float64

	AwayBookie   string
	AwayWinDelta float64
	AwayWinOdds  float64
}
