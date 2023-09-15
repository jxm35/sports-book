package predictions

import (
	"context"
	"fmt"
	"sync"

	"sports-book.com/pkg/config"
	"sports-book.com/pkg/domain"
)

var (
	predictionHandler PredictionHandler
	predictionOnce    sync.Once
)

type PredictionHandler interface {
	SavePrediction(ctx context.Context, matchId int32, prediction domain.Prediction) error
	GetPrediction(ctx context.Context, matchId int32) (domain.Prediction, error)
}

func GetPredictionHandler() PredictionHandler {
	predictionOnce.Do(
		func() {
			impl, found := config.GetConfigVal[string]("predictions.impl").Get()
			if !found {
				panic(config.ErrConfigNotProvided)
			}
			switch impl {
			case "in_memory":
				predictionHandler = newInMemoryPredictionHandler()
			case "db":
				predictionHandler = newDatabasePredictionHandler()
			default:
				panic(fmt.Sprintf("invalid prediction type %s", impl))
			}
		},
	)
	return predictionHandler
}
