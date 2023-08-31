package predict

import (
	"github.com/jxm35/go-results"
	"sports-book.com/predict/domain"
)

type Pipeline interface {
	PredictMatch(homeTeam, awayTeam, season int32) (domain.MatchProbability, OddsDelta, error)
	PlaceBet(matchId int32, generatedOdds domain.MatchProbability, currentPot float64) results.Option[domain.BetOrder]
}
