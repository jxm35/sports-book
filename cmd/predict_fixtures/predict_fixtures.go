package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"sports-book.com/pkg/bet_selector"
	"sports-book.com/pkg/db"
	"sports-book.com/pkg/domain"
	"sports-book.com/pkg/notify"
	"sports-book.com/pkg/pipeline"
	"sports-book.com/pkg/probability_generator"
	"sports-book.com/pkg/score_predictor"
)

func main() {
	_, err := db.Connect()
	if err != nil {
		panic(err)
	}

	predictionPipeline, err := pipeline.NewPipelineBuilder().
		SetPredictor(score_predictor.NewEloGoalsPredictor(5, 11)).
		// SetPredictor(&goals_predictor.LastSeasonXgGoalPredictor{LastXGames: 0}).
		SetProbabilityGenerator(&probability_generator.WeibullOddsGenerator{}).
		SetBetPlacer(bet_selector.NewKellyCriterionBetSelector(0.1, 0.3, 0.05, true)).
		// SetBetPlacer(bet_selector.NewFixedAmountBetSelector(0.1, 0.3, 0.2)).
		Build()
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

		if err := json.Unmarshal([]byte(record.Body), &fixtures); err != nil {
			fmt.Println("Failed to unmarshal fixtures: ", record.Body)
			return fmt.Errorf("failed to unmarshal fixtures: %w", err)
		}
		for _, fixture := range fixtures {
			match, err := db.CreateFixture(ctx, fixture, -1)
			if err != nil {
				return fmt.Errorf("failed to add fixture: %w", err)
			}
			probabilities, _, err := predictionPipeline.PredictMatch(
				match.HomeTeam,
				match.AwayTeam,
				2023,
				domain.LeagueEPL,
				time.Now(),
				match.ID,
			)
			bet := predictionPipeline.PlaceBet(match.ID, probabilities, 100)
			if bet.IsPresent() {
				if err := notify.GetNotifier().NotifyBetPlaced(bet.Value()); err != nil {
					return fmt.Errorf("failed to add notify bet order: %w", err)
				}
			}
		}
	}

	return nil
}
