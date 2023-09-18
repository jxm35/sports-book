package new_results

import (
	"context"
	"fmt"
	"strconv"

	"sports-book.com/pkg/db"
	"sports-book.com/pkg/domain"
	"sports-book.com/pkg/gorm/model"
	"sports-book.com/pkg/logger"
)

func HandleNewResults(ctx context.Context, results []domain.Result) error {
	for _, result := range results {
		// update fixture to be a result
		usId, err := strconv.Atoi(result.Id)
		if err != nil {
			return err
		}
		match, err := db.GetMatchByUsId(ctx, int32(usId))
		if match.HomeTeam == 0 || match.AwayTeam == 0 {
			logger.Error("match not found in database")
			return fmt.Errorf("match not found in database")
		}
		homeGoals, err := strconv.Atoi(result.Goals.Home)
		if err != nil {
			logger.Error("couldn't get homeGoals from message", "error", err)
			return err
		}
		awayGoals, err := strconv.Atoi(result.Goals.Away)
		if err != nil {
			fmt.Println("couldn't get awayGoals from message", "error", err)
			return err
		}
		homeXg, err := strconv.ParseFloat(result.XG.Home, 64)
		if err != nil {
			logger.Error("couldn't get home xG from message", "error", err)
			return err
		}
		awayXg, err := strconv.ParseFloat(result.XG.Away, 64)
		if err != nil {
			logger.Error("couldn't get away xG from message", "error", err)
			return err
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
			logger.Error("couldn't update match", "match", match.ID, "error", err)
			return err
		}

		// calculate the returns of the bet if any / update the pot
	}
	return nil
}
