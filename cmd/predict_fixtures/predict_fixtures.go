package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"sports-book.com/pkg/db"
	"sports-book.com/pkg/domain"
	"sports-book.com/pkg/notify"
	"sports-book.com/pkg/pipeline"
)

func main() {
	_, err := db.Connect()
	if err != nil {
		panic(err)
	}

	predictionPipeline, err := pipeline.NewPipelineFromConfig()
	if err != nil {
		panic(err)
	}

	lambda.Start(func(ctx context.Context, event events.SQSEvent) error {
		return handle(ctx, event, predictionPipeline)
	})
}

func handle(ctx context.Context, event events.SQSEvent, predictionPipeline pipeline.Pipeline) error {
	for _, record := range event.Records {
		var fixtures []domain.Fixture
		bets := make([]domain.BetOrder, 0)

		if err := json.Unmarshal([]byte(record.Body), &fixtures); err != nil {
			fmt.Println("Failed to unmarshal fixtures: ", record.Body)
			return fmt.Errorf("failed to unmarshal fixtures: %w", err)
		}

		for _, fixture := range fixtures {
			match, err := db.CreateFixture(ctx, fixture, -1)
			if err != nil {
				return fmt.Errorf("failed to add fixture: %w", err)
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
			bet := predictionPipeline.PlaceBet(ctx, match.ID, probabilities, 100)
			if bet.IsPresent() {
				bets = append(bets, bet.Value())
			}
		}
		// sleep to ensure all odds have been entered into the database
		time.Sleep(5 * time.Second)

		for _, bet := range bets {
			if err := notify.GetNotifier().NotifyBetPlaced(ctx, bet); err != nil {
				return fmt.Errorf("failed to add notify bet order: %w", err)
			}

			// enter bet into the database
			if err := db.SaveBetPlaced(ctx, bet); err != nil {
				return fmt.Errorf("failed to save placed bet: %w", err)
			}
		}
	}

	return nil
}
