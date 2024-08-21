package dispatcher

import (
	"errors"
	tp "github.com/dsrvlabs/vatz/types"
	"strings"
	"sync"
	"time"

	pb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/dsrvlabs/vatz/manager/config"
	"github.com/dsrvlabs/vatz/utils"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
)

/* TODO: Discussion.
We need to discuss about notificatino module.
As I see this code, dispatcher itself is described is dispatcher
but dispatcher and dispatcher module should be splitted into two part.
*/

const (
	emojiER       string = "🚨"
	emojiDoubleEX string = "‼️"
	emojiSingleEx string = "❗"
	emojiCheck    string = "✅"
)

var (
	dispatcherSingletons []Dispatcher
	dispatcherOnce       sync.Once
)

// Dispatcher Notification provides interfaces to send alert dispatcher message with variable channel.
type Dispatcher interface {
	SetDispatcher(firstExecution bool, notificationFlag string, previousFlag tp.StateFlag, notifyInfo tp.NotifyInfo) error
	SendNotification(request tp.ReqMsg) error
}

// GetDispatchers gets the registered alert channel.
func GetDispatchers(cfg config.NotificationInfo) []Dispatcher {
	if len(cfg.DispatchChannels) == 0 {
		dpError := errors.New("error: No Dispatcher has set")
		log.Error().Str("module", "dispatcher").Msg("Please, Set at least a single channel for dispatcher, e.g.) Discord or Telegram")
		panic(dpError)
	}

	dispatcherOnce.Do(func() {
		for _, chainInfo := range cfg.DispatchChannels {
			var chainNotificationFlag = ""
			if len(chainInfo.ReminderSchedule) == 0 {
				chainInfo.ReminderSchedule = cfg.DefaultReminderSchedule
			}

			if chainInfo.Flag != "" {
				log.Debug().Str("module", "dispatcher").Msgf("plugin Flag exists!: %s", chainInfo.Flag)
				chainNotificationFlag = chainInfo.Flag
			}
			switch channel := chainInfo.Channel; {
			case strings.EqualFold(channel, string(tp.Discord)):
				dispatcherSingletons = append(dispatcherSingletons, &discord{
					host:             cfg.HostName,
					channel:          tp.Discord,
					secret:           chainInfo.Secret,
					notificationFlag: chainNotificationFlag,
					reminderCron:     cron.New(cron.WithLocation(time.UTC)),
					reminderSchedule: chainInfo.ReminderSchedule,
					entry:            sync.Map{},
				})
			case strings.EqualFold(channel, string(tp.Telegram)):
				dispatcherSingletons = append(dispatcherSingletons, &telegram{
					host:             cfg.HostName,
					channel:          tp.Telegram,
					secret:           chainInfo.Secret,
					chatID:           chainInfo.ChatID,
					notificationFlag: chainNotificationFlag,
					reminderCron:     cron.New(cron.WithLocation(time.UTC)),
					reminderSchedule: chainInfo.ReminderSchedule,
					entry:            sync.Map{},
				})
			case strings.EqualFold(channel, string(tp.Slack)):
				dispatcherSingletons = append(dispatcherSingletons, &slack{
					host:             cfg.HostName,
					channel:          tp.Slack,
					secret:           chainInfo.Secret,
					notificationFlag: chainNotificationFlag,
					reminderCron:     cron.New(cron.WithLocation(time.UTC)),
					reminderSchedule: chainInfo.ReminderSchedule,
					entry:            sync.Map{},
				})
			case strings.EqualFold(channel, string(tp.PagerDuty)):
				dispatcherSingletons = append(dispatcherSingletons, &pagerduty{
					host:             cfg.HostName,
					channel:          tp.PagerDuty,
					secret:           chainInfo.Secret,
					notificationFlag: chainNotificationFlag,
					pagerEntry:       sync.Map{},
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

	pUnique := utils.MakeUniqueValue(info.Plugin, info.Address, info.Port)

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
		Options:      map[string]interface{}{"pUnique": pUnique},
	}
}
