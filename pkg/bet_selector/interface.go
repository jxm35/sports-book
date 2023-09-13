package bet_selector

import (
	results "github.com/jxm35/go-results"

	"sports-book.com/pkg/domain"
)

type BetSelector interface {
	Place1x2Bets(matchId int32, generatedOdds domain.MatchProbability, currentPot float64) results.Option[domain.BetOrder]
}
