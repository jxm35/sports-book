package domain

type BackingType string

var (
	BackHomeWin BackingType = "BackHomeWin"
	BackDraw    BackingType = "BackDraw"
	BackAwayWin BackingType = "BackAwayWin"
)

const (
	LeagueEPL        = "epl"
	LeagueLaLiga     = "la_liga"
	LeagueBundesliga = "bundesliga"
	LeagueSerieA     = "serie_a"
)

type MatchProbability struct {
	HomeWin float64
	Draw    float64
	AwayWin float64
}

type BetOrder struct {
	MatchId   int32
	Backing   BackingType
	BookMaker string
	OddsTaken float64
	Amount    float64
}
