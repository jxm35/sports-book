package predict

import (
	"time"

	results "github.com/jxm35/go-results"

	"sports-book.com/predict/domain"
)

type Pipeline interface {
	PredictMatch(homeTeam, awayTeam, season int32, league string, date time.Time, matchID int32) (domain.MatchProbability, OddsDelta, error)
	PlaceBet(matchId int32, generatedOdds domain.MatchProbability, currentPot float64) results.Option[domain.BetOrder]
}
