package pipeline

import (
	"time"

	results "github.com/jxm35/go-results"

	"sports-book.com/pkg/domain"
)

type Pipeline interface {
	PredictMatch(homeTeam, awayTeam, season int32, league domain.League, date time.Time, matchID int32) (domain.MatchProbability, domain.OddsDelta, error)
	PlaceBet(matchId int32, generatedOdds domain.MatchProbability, currentPot float64) results.Option[domain.BetOrder]
}
