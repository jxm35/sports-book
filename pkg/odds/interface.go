package odds

import (
	"context"
	"sync"

	"sports-book.com/pkg/config"
	"sports-book.com/pkg/domain"
)

var (
	oddsRetriever     OddsRetriever
	oddsRetrieverOnce sync.Once
)

type OddsRetriever interface {
	GetBestOdds(ctx context.Context, matchId int32, league domain.League) (domain.BookmakerMatchOdds, error)
}

func GetOddsRetriever() OddsRetriever {
	oddsRetrieverOnce.Do(
		func() {
			impl, found := config.GetConfigVal[string]("odds_retriever.impl").Get()
			if !found {
				panic(config.ErrConfigNotProvided)
			}
			switch impl {
			case "db":
				oddsRetriever = &dbRetriever{}
			case "api":
				oddsRetriever = &apiRetriever{}
			default:
				panic("invalid odds retriever config")
			}
		},
	)
	return oddsRetriever
}
