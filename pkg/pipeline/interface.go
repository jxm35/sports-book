package pipeline

import (
	"context"
	"time"

	results "github.com/jxm35/go-results"

	"sports-book.com/pkg/domain"
)

type Pipeline interface {
	PredictMatch(ctx context.Context, homeTeam, awayTeam, season int32, league domain.League, date time.Time, matchID int32) (domain.MatchProbability, error)
	PlaceBet(ctx context.Context, matchId int32, generatedOdds domain.MatchProbability, currentPot float64) results.Option[domain.BetOrder]
}
