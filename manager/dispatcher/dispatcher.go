package dispatcher

import (
	"github.com/dsrvlabs/vatz/manager/config"
	tp "github.com/dsrvlabs/vatz/manager/types"
	"sync"
)

/* TODO: Discussion.
We need to discuss about notificatino module.
As I see this code, dispatcher itself is described is dispatcher
but dispatcher and dispatcher module should be splitted into two part.
*/

var (
	dispatcherSingleton Dispatcher
	dispatcherOnce      sync.Once
)

// Dispatcher Notification provides interfaces to send alert dispatcher message with variable channel.
type Dispatcher interface {
	SendNotification(request tp.ReqMsg) error
}

func GetDispatchers(cfg config.NotificationInfo) []Dispatcher {
	// !!Note!!
	// This is a sample and will be modified with issue #226
	// Please, remove this comment when you create a channel with notification.
	sample1 := &discord{channel: tp.Discord}

	type sampleSecret struct {
		secret  string
		channel tp.Channel
	}

	sampleSecrets := []sampleSecret{
		{cfg.DiscordSecret, tp.Discord},
	}

	var dispatchers []Dispatcher
	dispatcherOnce.Do(func() {
		for _, secretInfo := range sampleSecrets {
			if secretInfo.channel == tp.Discord {
				dispatcherSingleton = sample1
				dispatchers = append(dispatchers, dispatcherSingleton)
			}
		}
	})
	return dispatchers
}
