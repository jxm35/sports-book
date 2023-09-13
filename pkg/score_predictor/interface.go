package score_predictor

import (
	"time"

	"sports-book.com/pkg/domain"
)

type ScorePredictor interface {
	PredictScore(homeTeam, awayTeam, season int32, league domain.League, date time.Time, matchID int32) (float64, float64, error)
}
