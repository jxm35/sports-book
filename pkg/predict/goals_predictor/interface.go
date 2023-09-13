package goals_predictor

import "time"

type GoalsPredictor interface {
	PredictScore(homeTeam, awayTeam, season int32, league string, date time.Time, matchID int32) (float64, float64, error)
}
