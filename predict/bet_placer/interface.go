package bet_placer

import (
	"github.com/jxm35/go-results"
	"sports-book.com/predict/domain"
)

type BetPlacer interface {
	Place1x2Bets(matchId int32, generatedOdds domain.MatchProbability, currentPot float64) results.Option[domain.BetOrder]
}
