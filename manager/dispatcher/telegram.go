package dispatcher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	pb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	tp "github.com/dsrvlabs/vatz/manager/types"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
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
	pUnique := deliverMessage.Options["pUnique"].(string)

	if reqToNotify {
		err := t.SendNotification(deliverMessage)
		if err != nil {
			log.Error().Str("module", "dispatcher").Msgf("Channel(Telegram): Send notification error: %s", err)
			return err
		}

	}

	if reminderState == tp.ON {
		newEntries := []cron.EntryID{}
		/*
			In case of reminder has to keep but stateFlag has changed,
			e.g.) CRITICAL -> WARNING
			e.g.) ERROR -> INFO -> ERROR
		*/
		if entries, ok := t.entry.Load(pUnique); ok {
			for _, entry := range entries.([]cron.EntryID) {
				t.reminderCron.Remove(entry)
			}
			t.reminderCron.Stop()
		}
		for _, schedule := range t.reminderSchedule {
			id, _ := t.reminderCron.AddFunc(schedule, func() {
				err := t.SendNotification(deliverMessage)
				if err != nil {
					log.Error().Str("module", "dispatcher").Msgf("Channel(Telegram): Send notification error: %s", err)
				}
			})
			newEntries = append(newEntries, id)
		}
		t.entry.Store(pUnique, newEntries)
		t.reminderCron.Start()
	} else if reminderState == tp.OFF {
		entries, _ := t.entry.Load(pUnique)
		for _, entity := range entries.([]cron.EntryID) {
			t.reminderCron.Remove(entity)
		}
		t.reminderCron.Stop()
	}
	return nil
}

func (t *telegram) SendNotification(msg tp.ReqMsg) error {
	var (
		err      error
		response *http.Response
		emoji    = emojiER
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

	url := fmt.Sprintf("%s/sendMessage", getURL(t.secret))
	sendingText := fmt.Sprintf(`
%s<strong>%s</strong>%s
<strong>(%s)</strong>
Plugin Name: <em>%s</em>
%s`, emoji, msg.Severity.String(), emoji, t.host, msg.ResourceType, msg.Msg)

	body, _ := json.Marshal(map[string]string{
		"chat_id":    t.chatID,
		"text":       sendingText,
		"parse_mode": "html",
	})

	response, err = http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Error().Str("module", "dispatcher").Msgf("dispatcher telegram Error: %s", err)
		return err
	}
	defer response.Body.Close()
	body, err = io.ReadAll(response.Body)
	if err != nil {
		log.Error().Str("module", "dispatcher").Msgf("Channel(Telegram): body parsing Error: %s", err)
		return err
	}
	respJSON := make(map[string]interface{})
	if !respJSON["ok"].(bool) {
		log.Error().Str("module", "dispatcher").Msg("Channel(Telegram): Connection failed due to Invalid telegram token.")
	}
	err = json.Unmarshal(body, &respJSON)
	if err != nil {
		log.Error().Str("module", "dispatcher").Msgf("Channel(Telegram): Unmarshalling JSON Error: %s", err)
		return err
	}
	return nil
}

func getURL(token string) string {
	return fmt.Sprintf("https://api.telegram.org/bot%s", token)
}
