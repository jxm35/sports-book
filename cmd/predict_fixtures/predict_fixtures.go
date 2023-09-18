package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"sports-book.com/internal/new_fixtures"
	"sports-book.com/internal/new_results"
	"sports-book.com/pkg/db"
	"sports-book.com/pkg/domain"
	"sports-book.com/pkg/logger"
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
		switch *title.StringValue {
		case "new_fixtures":
			var fixtures []domain.Fixture

			if err := json.Unmarshal([]byte(record.Body), &fixtures); err != nil {
				logger.Error("failed to unmarshal fixtures", "body", record.Body)
				return err
			}

			return new_fixtures.HandleNewFixtures(ctx, fixtures, predictionPipeline)

		case "new_results":
			var results []domain.Result

			if err := json.Unmarshal([]byte(record.Body), &results); err != nil {
				logger.Error("failed to unmarshal results", "body", record.Body)
				return err
			}

			return new_results.HandleNewResults(ctx, results)

		default:
			logger.Error("unknown event type", "title", *title.StringValue)
			return fmt.Errorf("unknown event type %s", *title.StringValue)

		}
	}
	return nil
}
