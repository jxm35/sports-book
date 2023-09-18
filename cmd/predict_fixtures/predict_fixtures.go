package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"sports-book.com/pkg/db"
	"sports-book.com/pkg/domain"
	"sports-book.com/pkg/logger"
	"sports-book.com/pkg/notify"
	"sports-book.com/pkg/pipeline"
)

func main() {
	logger.InitialiseStructuredLogger()
	logger.Info("logs started")

	_, err := db.Connect()
	if err != nil {
		panic(err)
	}

	predictionPipeline, err := pipeline.NewPipelineFromConfig()
	if err != nil {
		panic(err)
	}

	lambda.Start(func(ctx context.Context, event events.SQSEvent) error {
		logger.Info("event received", "event", event)
		return handle(ctx, event, predictionPipeline)
	})
}

func handle(ctx context.Context, event events.SQSEvent, predictionPipeline pipeline.Pipeline) error {
	for _, record := range event.Records {
		title := record.MessageAttributes["Title"]
		logger.Info("new message received", "title", *title.StringValue)
		var fixtures []domain.Fixture
		bets := make([]domain.BetOrder, 0)

		if err := json.Unmarshal([]byte(record.Body), &fixtures); err != nil {
			logger.Error("failed to unmarshal fixtures", "body", record.Body)
			return err
		}

		for _, fixture := range fixtures {
			match, err := db.CreateFixture(ctx, fixture, -1)
			if err != nil {
				logger.Error("failed to save fixture", "error", err)
				return err
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
			if err != nil {
				logger.Error("failed to predict match", "error", err)
				return err
			}

			bet := predictionPipeline.PlaceBet(ctx, match.ID, probabilities, 100)
			if bet.IsPresent() {
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
	}

	return nil
}
