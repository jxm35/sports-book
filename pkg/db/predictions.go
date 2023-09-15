package db

import (
	"context"

	"sports-book.com/pkg/domain"
	"sports-book.com/pkg/gorm/model"
	"sports-book.com/pkg/gorm/query"
)

func GetPrediction(ctx context.Context, matchId int32) (model.Prediction, error) {
	p := query.Prediction
	res, err := p.WithContext(ctx).First()
	if err != nil {
		return model.Prediction{}, err
	}
	return *res, nil
}

func SavePrediction(ctx context.Context, matchId int32, prediction domain.Prediction) error {
	p := query.Prediction
	err := p.WithContext(ctx).
		Create(&model.Prediction{
			MatchID: matchId,
			HomeXg:  prediction.HomexG,
			AwayXg:  prediction.AwayxG,
		})
	return err
}
