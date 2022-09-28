package dispatcher

import (
	"fmt"
	"sync"

	"github.com/dsrvlabs/vatz/manager/config"
	tp "github.com/dsrvlabs/vatz/manager/types"
)

/* TODO: Discussion.
We need to discuss about notificatino module.
As I see this code, dispatcher itself is described is dispatcher
but dispatcher and dispatcher module should be splitted into two part.
*/

var (
	dispatcherSingletons []Dispatcher
	dispatcherOnce       sync.Once
)

// Dispatcher Notification provides interfaces to send alert dispatcher message with variable channel.
type Dispatcher interface {
	SendNotification(request tp.ReqMsg) error
}

func GetDispatchers(cfg config.NotificationInfo) []Dispatcher {
	dispatcherOnce.Do(func() {
		for _, chanInfo := range cfg.DispatchChannels {
			switch chanInfo.Channel {
			case "discord":
				discord := &discord{
					channel: tp.Discord,
					secret:  chanInfo.Secret,
				}
				dispatcherSingletons = append(dispatcherSingletons, discord)
			case "telegram":
				telegram := &telegram{
					channel: tp.Telegram,
					secret:  chanInfo.Secret,
					chatID:  chanInfo.ChatID,
				}
				dispatcherSingletons = append(dispatcherSingletons, telegram)
			default:
				fmt.Println(chanInfo.Channel, "is not work")
			}
		}
	})
	return dispatcherSingletons
}
