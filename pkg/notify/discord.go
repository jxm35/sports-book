package notify

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/gtuk/discordwebhook"

	"sports-book.com/pkg/db"
	"sports-book.com/pkg/domain"
)

var ErrNoUrlFound = errors.New("no webhook url found for discord bot")

type discordNotifier struct {
	username string
	url      string
}

func newDiscordNotifier() (*discordNotifier, error) {
	url := os.Getenv("DISCORD_WEBHOOK_URL")
	if url == "" {
		return nil, ErrNoUrlFound
	}
	return &discordNotifier{
		username: "sports-book-bot",
		url:      url,
	}, nil
}

func (d *discordNotifier) NotifyBetPlaced(ctx context.Context, bet domain.BetOrder) error {
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
	return d.sendMessage(message)
}

func (l *discordNotifier) NotifyError(message string) error {
	return nil
}

func (d *discordNotifier) sendMessage(text string) error {
	user := d.username

	message := discordwebhook.Message{
		Username: &user,
		Content:  &text,
	}

	if err := discordwebhook.SendMessage(d.url, message); err != nil {
		return err
	}
	return nil
}
