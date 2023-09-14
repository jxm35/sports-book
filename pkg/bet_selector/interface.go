package bet_selector

import (
	"sync"

	results "github.com/jxm35/go-results"

	"sports-book.com/pkg/config"
	"sports-book.com/pkg/domain"
)

var (
	betSelector     BetSelector
	betSelectorOnce sync.Once
)

type BetSelector interface {
	Place1x2Bets(matchId int32, generatedOdds domain.MatchProbability, currentPot float64) results.Option[domain.BetOrder]
}

func NewBetSelectorFromConfig() (BetSelector, error) {
	var err error
	betSelectorOnce.Do(
		func() {
			impl := config.GetConfigVal[string]("bet_selector.impl")
			if impl.IsNone() {
				err = config.ErrConfigNotProvided
				return
			}
			switch impl.Value() {
			case "fixed_amount":
				selector, e := NewFixedAmountBetSelectorFromConfig()
				err = e
				betSelector = selector
			case "kelly_criterion":
				selector, e := NewKellyCriterionBetSelectorFromConfig()
				err = e
				betSelector = selector
			default:
				err = config.ErrConfigNotProvided
			}
		})
	if err != nil {
		return nil, err
	}
	return betSelector, nil
}
