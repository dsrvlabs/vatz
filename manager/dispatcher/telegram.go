package dispatcher

import (
	"bytes"
	"encoding/json"
	"fmt"
	pb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	tp "github.com/dsrvlabs/vatz/manager/types"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"sync"
)

// telegram: This is a sample code
// that helps to multi methods for notification.
type telegram struct {
	host             string
	channel          tp.Channel
	secret           string
	chatID           string
	reminderSchedule []string
	reminderCron     *cron.Cron
	entry            sync.Map
}

func (t *telegram) SetDispatcher(firstRunMsg bool, preStat tp.StateFlag, notifyInfo tp.NotifyInfo) error {
	reqToNotify, reminderState, deliverMessage := messageHandler(firstRunMsg, preStat, notifyInfo)

	if reqToNotify {
		t.SendNotification(deliverMessage)
	}

	if reminderState == tp.ON {
		newEntries := []cron.EntryID{}
		//In case of reminder has to keep but stateFlag has changed,
		//e.g.) CRITICAL -> WARNING
		//e.g.) ERROR -> INFO -> ERROR
		if entries, ok := t.entry.Load(notifyInfo.Method); ok {
			for _, entry := range entries.([]cron.EntryID) {
				t.reminderCron.Remove(entry)
			}
			t.reminderCron.Stop()
		}
		for _, schedule := range t.reminderSchedule {
			id, _ := t.reminderCron.AddFunc(schedule, func() {
				t.SendNotification(deliverMessage)
			})
			newEntries = append(newEntries, id)
		}
		t.entry.Store(notifyInfo.Method, newEntries)
		t.reminderCron.Start()

	} else if reminderState == tp.OFF {
		entries, _ := t.entry.Load(notifyInfo.Method)
		for _, entity := range entries.([]cron.EntryID) {
			t.reminderCron.Remove(entity)
		}
		t.reminderCron.Stop()
	}
	return nil
}

func (t *telegram) SendNotification(msg tp.ReqMsg) error {
	var err error
	var response *http.Response
	emoji := "🚨"
	if msg.State == pb.STATE_SUCCESS {
		if msg.Severity == pb.SEVERITY_CRITICAL {
			emoji = "‼️"
		} else if msg.Severity == pb.SEVERITY_WARNING {
			emoji = "❗"
		} else if msg.Severity == pb.SEVERITY_INFO {
			emoji = "✅"
		}
	}
	url := fmt.Sprintf("%s/sendMessage", getUrl(t.secret))
	sendingText := fmt.Sprintf(`
%s**%s**%s
**(%s)**
_Plugin Name: %s_
%s`, emoji, msg.Severity.String(), emoji, t.host, msg.ResourceType, msg.Msg)

	body, _ := json.Marshal(map[string]string{
		"chat_id":    t.chatID,
		"text":       sendingText,
		"parse_mode": "markdown",
	})

	response, err = http.Post(
		url,
		"application/json",
		bytes.NewBuffer(body),
	)

	if err != nil {
		log.Error().Str("module", "dispatcher").Msgf("dispatcher telegram Error: %s", err)
		return err
	}
	defer response.Body.Close()

	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error().Str("module", "dispatcher").Msgf("dispatcher telegram body parsing Error: %s", err)
		return err
	} else {
		respJSON := make(map[string]interface{})
		json.Unmarshal(body, &respJSON)
		if !respJSON["ok"].(bool) {
			log.Error().Str("module", "dispatcher").Msg("dispatcher CH: Telegram-Invalid telegram token.")
		}
	}
	return nil
}

func getUrl(token string) string {
	return fmt.Sprintf("https://api.telegram.org/bot%s", token)
}
