package new_fixtures

import (
	"context"
	"errors"
	"strconv"
	"time"

	"sports-book.com/pkg/db"
	"sports-book.com/pkg/domain"
	"sports-book.com/pkg/logger"
	"sports-book.com/pkg/notify"
	"sports-book.com/pkg/pipeline"
	"sports-book.com/pkg/score_predictor"
)

func HandleNewFixtures(ctx context.Context, fixtures []domain.Fixture, predictionPipeline pipeline.Pipeline) error {
	bets := make([]domain.BetOrder, 0)
	for _, fixture := range fixtures {
		usId, err := strconv.Atoi(fixture.Id)
		if err != nil {
			logger.Error("could not retrieve understat id", "error", err)
			return err
		}
		match, err := db.GetMatchByUsId(ctx, int32(usId))
		if err != nil {
			match, err = db.CreateFixture(ctx, fixture, 38)
			if err != nil {
				logger.Error("failed to save fixture", "error", err)
				return err
			}
		}
		probabilities, err := predictionPipeline.PredictMatch(
			ctx,
			match.HomeTeam,
			match.AwayTeam,
			2023,
			domain.LeagueEPL,
			time.Now(),
			match.ID,
		)
		if errors.Is(err, score_predictor.ErrNoPreviousData) {
			logger.Debug("no data for one or both teams, continuing", "match_id", match.ID)
			continue
		} else if err != nil {
			logger.Error("failed to predict match", "error", err)
			return err
		}
		logger.Debug("match predicted", "match_id", match.ID, "probabilities", probabilities)

		bet := predictionPipeline.PlaceBet(ctx, match.ID, probabilities, 100)
		if bet.IsPresent() {
			logger.Debug("bet placed", "match", match.ID, "bet", bet.Value())
			bets = append(bets, bet.Value())
		}
	}
	// sleep to ensure all odds have been entered into the database
	time.Sleep(5 * time.Second)

	for _, bet := range bets {
		// send message to me recommending for bet to be placed
		if err := notify.GetNotifier().NotifyBetPlaced(ctx, bet); err != nil {
			logger.Error("failed to notify about placed bet", "error", err)
			return err
		}

		// enter bet into the database
		if err := db.SaveBetPlaced(ctx, bet); err != nil {
			logger.Error("could not save bet", "error", err)
			return err
		}
	}
	return nil
}
