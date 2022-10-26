package dispatcher

import (
	"errors"
	"strings"
	"sync"
	"time"

	pb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/dsrvlabs/vatz/manager/config"
	tp "github.com/dsrvlabs/vatz/manager/types"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
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
	SetDispatcher(firstExecution bool, port int, previousFlag tp.StateFlag, notifyInfo tp.NotifyInfo) error
	SendNotification(request tp.ReqMsg) error
}

func GetDispatchers(cfg config.NotificationInfo) []Dispatcher {
	if len(cfg.DispatchChannels) == 0 {
		dpError := errors.New("Error: No Dispatcher has set.")
		log.Error().Str("module", "dispatcher").Msg("Please, Set at least a channel for dispatcher, e.g.) Discord or Telegram")
		panic(dpError)
	}

	dispatcherOnce.Do(func() {
		for _, chanInfo := range cfg.DispatchChannels {
			if len(chanInfo.ReminderSchedule) == 0 {
				chanInfo.ReminderSchedule = cfg.DefaultReminderSchedule
			}
			switch channel := chanInfo.Channel; {
			case strings.EqualFold(channel, string(tp.Discord)):
				dispatcherSingletons = append(dispatcherSingletons, &discord{
					host:             cfg.HostName,
					channel:          tp.Discord,
					secret:           chanInfo.Secret,
					reminderCron:     cron.New(cron.WithLocation(time.UTC)),
					reminderSchedule: chanInfo.ReminderSchedule,
					entry:            sync.Map{},
				})
			case strings.EqualFold(channel, string(tp.Telegram)):
				dispatcherSingletons = append(dispatcherSingletons, &telegram{
					host:             cfg.HostName,
					channel:          tp.Telegram,
					secret:           chanInfo.Secret,
					chatID:           chanInfo.ChatID,
					reminderCron:     cron.New(cron.WithLocation(time.UTC)),
					reminderSchedule: chanInfo.ReminderSchedule,
					entry:            sync.Map{},
				})
			}
		}
	})

	return dispatcherSingletons
}

func messageHandler(isFirst bool, preStat tp.StateFlag, info tp.NotifyInfo) (bool, tp.Reminder, tp.ReqMsg) {
	notifyOn := false
	reminderState := tp.HANG
	isFlagStateChanged := false

	if preStat.State != info.State || preStat.Severity != info.Severity {
		isFlagStateChanged = true
	}

	if info.State == pb.STATE_FAILURE ||
		(info.State == pb.STATE_SUCCESS && info.Severity == pb.SEVERITY_WARNING) ||
		(info.State == pb.STATE_SUCCESS && info.Severity == pb.SEVERITY_CRITICAL) {
		if isFirst || isFlagStateChanged {
			notifyOn = true
			reminderState = tp.ON
		}
	} else if info.State == pb.STATE_SUCCESS && info.Severity == pb.SEVERITY_INFO && isFlagStateChanged {
		notifyOn = true
		reminderState = tp.OFF
	}

	return notifyOn, reminderState, tp.ReqMsg{
		FuncName:     info.Method,
		State:        info.State,
		Msg:          info.ExecuteMsg,
		Severity:     info.Severity,
		ResourceType: info.Plugin,
	}
}
