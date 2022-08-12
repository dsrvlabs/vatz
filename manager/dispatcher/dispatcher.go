package dispatcher

import (
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

type dispatcher struct {
}

// GetDispatcher create new dispatcher dispatcher.
func GetDispatcher() Dispatcher {
	dispatcherOnce.Do(func() {
		dispatcherSingleton = &dispatcher{}
	})

	return dispatcherSingleton
}
