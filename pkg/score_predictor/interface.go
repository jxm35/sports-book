package score_predictor

import (
	"context"
	"sync"
	"time"

	"sports-book.com/pkg/config"
	"sports-book.com/pkg/domain"
)

var (
	scorePredictor     ScorePredictor
	scorePredictorOnce sync.Once
)

type ScorePredictor interface {
	PredictScore(ctx context.Context, homeTeam, awayTeam, season int32, league domain.League, date time.Time, matchID int32) (float64, float64, error)
}

func NewScorePredictorFromConfig() (ScorePredictor, error) {
	var err error
	scorePredictorOnce.Do(
		func() {
			impl := config.GetConfigVal[string]("score_predictor.impl")
			if impl.IsNone() {
				err = config.ErrConfigNotProvided
				return
			}
			switch impl.Value() {
			case "last_season_xg":
				scorePredictor = &LastSeasonXgScorePredictor{LastXGames: 0}
			case "elo":
				scorePredictor = NewEloGoalsPredictor(5, 11)
			default:
				err = config.ErrConfigNotProvided
			}
		})
	if err != nil {
		return nil, err
	}
	return scorePredictor, nil
}
