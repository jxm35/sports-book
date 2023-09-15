package notify

import (
	"context"
	"fmt"
	"log"

	"sports-book.com/pkg/db"
	"sports-book.com/pkg/domain"
)

type logNotifier struct{}

func (l *logNotifier) NotifyBetPlaced(ctx context.Context, bet domain.BetOrder) error {
	match, err := db.GetMatch(ctx, bet.MatchId)
	if err != nil {
		return err
	}
	homeTeam, err := db.GetTeam(ctx, match.HomeTeam)
	if err != nil {
		return err
	}
	awayTeam, err := db.GetTeam(ctx, match.AwayTeam)
	if err != nil {
		return err
	}
	message := fmt.Sprintf(`
Reccomended Bet:
%s
%s vs %s @ %s
Odds: %.2f from %s
Reccomended Stake: £%.2f
Predicted Probability: %.2f%%`,
		string(bet.Backing),
		homeTeam.Name, awayTeam.Name, match.Date,
		bet.OddsTaken, bet.BookMaker,
		bet.Amount,
		bet.PredictedProbability*100,
	)
	log.Print(message)
	return nil
}

func (l *logNotifier) NotifyError(message string) error {
	return nil
}
