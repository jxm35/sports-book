package goals_predictor

type GoalsPredictor interface {
	PredictScore(homeTeam, awayTeam, season int32, league string) (float64, float64, error)
}
