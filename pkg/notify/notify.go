package notify

import "sports-book.com/pkg/domain"

type notifier interface {
	NotifyBetPlaced(bet domain.BetOrder) error
}

func NotifyBetPlaced(bet domain.BetOrder) error {
	noti := &logNotifier{}
	return noti.NotifyBetPlaced(bet)
}
