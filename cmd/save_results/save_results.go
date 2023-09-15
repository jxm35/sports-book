package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"sports-book.com/pkg/db"
	"sports-book.com/pkg/domain"
	"sports-book.com/pkg/gorm/model"
)

func main() {
	_, err := db.Connect()
	if err != nil {
		panic(err)
	}

	lambda.Start(func(ctx context.Context, event events.SQSEvent) error {
		return handle(ctx, event)
	})
}

func handle(ctx context.Context, event events.SQSEvent) error {
	for _, record := range event.Records {
		var results []domain.Result

		if err := json.Unmarshal([]byte(record.Body), &results); err != nil {
			fmt.Println("Failed to unmarshal results: ", record.Body)
			return fmt.Errorf("failed to unmarshal results: %w", err)
		}
		for _, result := range results {
			// update fixture to be a result
			usId, err := strconv.Atoi(result.Id)
			if err != nil {
				fmt.Println("Failed to unmarshal results: ", record.Body)
				return fmt.Errorf("failed to unmarshal results: %w", err)
			}
			match, err := db.GetMatchByUsId(ctx, int32(usId))
			if match.HomeTeam == 0 || match.AwayTeam == 0 {
				fmt.Println("match not found in database")
				return fmt.Errorf("match not found in database")
			}
			homeGoals, err := strconv.Atoi(result.Goals.Home)
			if err != nil {
				fmt.Println("couldn't get homeGoals from message")
				return fmt.Errorf("couldn't get homeGoals from message: %w", err)
			}
			awayGoals, err := strconv.Atoi(result.Goals.Away)
			if err != nil {
				fmt.Println("couldn't get awayGoals from message")
				return fmt.Errorf("couldn't get awayGoals from message: %w", err)
			}
			homeXg, err := strconv.ParseFloat(result.XG.Home, 64)
			if err != nil {
				fmt.Println("couldn't get home xG from message")
				return fmt.Errorf("couldn't get home xG from message: %w", err)
			}
			awayXg, err := strconv.ParseFloat(result.XG.Away, 64)
			if err != nil {
				fmt.Println("couldn't get away xG from message")
				return fmt.Errorf("couldn't get away xG from message: %w", err)
			}
			err = db.UpdateMatch(ctx, match.ID, model.Match{
				ID:                match.ID,
				Date:              match.Date,
				HomeTeam:          match.HomeTeam,
				AwayTeam:          match.AwayTeam,
				Competition:       match.Competition,
				HomeGoals:         int32(homeGoals),
				AwayGoals:         int32(awayGoals),
				HomeExpectedGoals: homeXg,
				AwayExpectedGoals: awayXg,
				UsID:              match.UsID,
			})
			if err != nil {
				fmt.Println("couldn't update match: ", match.ID)
				return fmt.Errorf("couldn't update match	: %w", err)
			}

			// calculate the returns of the bet if any / update the pot
		}
	}

	return nil
}
