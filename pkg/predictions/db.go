package predictions

import (
	"context"

	"sports-book.com/pkg/db"
	"sports-book.com/pkg/domain"
)

type databasePredictionHandler struct{}

func newDatabasePredictionHandler() PredictionHandler {
	return &databasePredictionHandler{}
}

func (d *databasePredictionHandler) SavePrediction(ctx context.Context, matchId int32, prediction domain.Prediction) error {
	return db.SavePrediction(ctx, matchId, prediction)
}

func (d *databasePredictionHandler) GetPrediction(ctx context.Context, matchId int32) (domain.Prediction, error) {
	pred, err := db.GetPrediction(ctx, matchId)
	if err != nil {
		return domain.Prediction{}, err
	}
	return domain.Prediction{
		HomexG: pred.HomeXg,
		AwayxG: pred.AwayXg,
	}, nil
}
