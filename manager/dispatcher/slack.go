package dispatcher

import (
	"bytes"
	"encoding/json"
	"fmt"
	tp "github.com/dsrvlabs/vatz/types"
	"github.com/dsrvlabs/vatz/utils"
	"net/http"
	"sync"

	pb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
)

// slack: This is a sample code
// that helps to multi methods for notification.
type slack struct {
	host             string
	channel          tp.Channel
	secret           string
	notificationFlag string
	reminderSchedule []string
	reminderCron     *cron.Cron
	entry            sync.Map
}

type SlackRequestBody struct {
	Text string `json:"text"`
}

func (s *slack) SetDispatcher(firstRunMsg bool, pluginNotificationFlag string, preStat tp.StateFlag, notifyInfo tp.NotifyInfo) error {
	reqToNotify, reminderState, deliverMessage := messageHandler(firstRunMsg, preStat, notifyInfo)
	pUnique := deliverMessage.Options["pUnique"].(string)
	flagEnabled, sameFlagExists := utils.IsNotifiedEnabledAndSend(s.notificationFlag, pluginNotificationFlag)
	if !flagEnabled || flagEnabled && sameFlagExists {
		if reqToNotify {
			err := s.SendNotification(deliverMessage)
			if err != nil {
				log.Error().Str("module", "dispatcher").Msgf("Channel(Slack): Send notification error: %s", err)
				return err
			}

		}
	}

	if reminderState == tp.ON {
		newEntries := []cron.EntryID{}
		/*
			In case of reminder has to keep but stateFlag has changed,
			e.g.) CRITICAL -> WARNING
			e.g.) ERROR -> INFO -> ERROR
		*/
		if entries, ok := s.entry.Load(pUnique); ok {
			for _, entry := range entries.([]cron.EntryID) {
				s.reminderCron.Remove(entry)
			}
			s.reminderCron.Stop()
		}
		for _, schedule := range s.reminderSchedule {
			id, _ := s.reminderCron.AddFunc(schedule, func() {
				err := s.SendNotification(deliverMessage)
				if err != nil {
					log.Error().Str("module", "dispatcher").Msgf("Channel(Slack): Send notification error: %s", err)
				}
			})
			newEntries = append(newEntries, id)
		}
		s.entry.Store(pUnique, newEntries)
		s.reminderCron.Start()
	} else if reminderState == tp.OFF {
		entries, _ := s.entry.Load(pUnique)
		if _, ok := entries.([]cron.EntryID); ok {
			for _, entity := range entries.([]cron.EntryID) {
				s.reminderCron.Remove(entity)
			}
			s.reminderCron.Stop()
		}
	}
	return nil
}

func (s *slack) SendNotification(msg tp.ReqMsg) error {
	var (
		err   error
		emoji = emojiER
	)

	if msg.State == pb.STATE_SUCCESS {
		switch {
		case msg.Severity == pb.SEVERITY_CRITICAL:
			emoji = emojiDoubleEX
		case msg.Severity == pb.SEVERITY_WARNING:
			emoji = emojiSingleEx
		case msg.Severity == pb.SEVERITY_INFO:
			emoji = emojiCheck
		}
	}
	sendingText := fmt.Sprintf(`
%s *%s* %s
> 
Host: *%s*
Plugin Name: _%s_
%s`, emoji, msg.Severity.String(), emoji, s.host, msg.ResourceType, msg.Msg)
	slackBody, _ := json.Marshal(SlackRequestBody{Text: sendingText})
	req, err := http.NewRequest(http.MethodPost, s.secret, bytes.NewBuffer(slackBody))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return err
	}
	if buf.String() != "ok" {
		return fmt.Errorf("non-ok response returned from Slack")
	}

	return nil
}
