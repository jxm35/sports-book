package notify

import (
	"sync"

	"sports-book.com/pkg/config"
	"sports-book.com/pkg/domain"
)

var (
	notify     notifier
	notifyOnce sync.Once
)

type notifier interface {
	NotifyBetPlaced(bet domain.BetOrder) error
}

func GetNotifier() notifier {
	notifyOnce.Do(
		func() {
			impl := config.GetConfigVal[string]("notify.impl").Value()
			switch impl {
			case "discord":
				n, err := newDiscordNotifier()
				if err != nil {
					panic(err)
				}
				notify = n
			default:
				notify = &logNotifier{}
			}
		},
	)
	return notify
}
