package predictions

import (
	"context"
	"errors"
	"sync"

	"sports-book.com/pkg/domain"
)

var ErrPredictionNotFound = errors.New("prediction not found")

type inMemoryPredictionHandler struct {
	matches map[int32]domain.Prediction
	lock    sync.RWMutex
}

func newInMemoryPredictionHandler() PredictionHandler {
	return &inMemoryPredictionHandler{
		matches: make(map[int32]domain.Prediction),
	}
}

func (i *inMemoryPredictionHandler) SavePrediction(ctx context.Context, matchId int32, prediction domain.Prediction) error {
	i.lock.Lock()
	defer i.lock.Unlock()
	i.matches[matchId] = prediction
	return nil
}

func (i *inMemoryPredictionHandler) GetPrediction(ctx context.Context, matchId int32) (domain.Prediction, error) {
	i.lock.RLock()
	defer i.lock.RUnlock()
	prediction, ok := i.matches[matchId]
	if !ok {
		return domain.Prediction{}, ErrPredictionNotFound
	}
	return prediction, nil
}
