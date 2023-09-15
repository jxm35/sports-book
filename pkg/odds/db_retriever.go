package odds

import (
	"context"

	"sports-book.com/pkg/db"
	"sports-book.com/pkg/domain"
)

type dbRetriever struct{}

func (d *dbRetriever) GetBestOdds(ctx context.Context, matchId int32, league domain.League) (domain.BookmakerMatchOdds, error) {
	return db.GetBestOddsForMatch(ctx, matchId)
}
