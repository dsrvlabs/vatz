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
	"time"
)

// telegram: This is a sample code
// that helps to multi methods for notification.
type telegram struct {
	host             string
	channel          tp.Channel
	secret           string
	chatID           string
	reminderSchedule []string
	reminderState    sync.Map
}

func (t telegram) SetDispatcher(firstRunMsg bool, preStat tp.StateFlag, notifyInfo tp.NotifyInfo) error {
	reqToNotify, setReminders, deliverMessage := messageHandler(firstRunMsg, preStat, notifyInfo)
	methodName := notifyInfo.Method
	if reqToNotify {
		t.SendNotification(deliverMessage)
	}

	if setReminders.ReminderState == tp.ON {
		fmt.Println("TP ON: Is going to Start!! 1")
		c := cron.New(cron.WithLocation(time.UTC))
		if _, ok := t.reminderState.Load(methodName); ok {
			preCron, _ := t.reminderState.Load(methodName)
			c = preCron.(*cron.Cron)
		}
		for _, schedule := range t.reminderSchedule {
			id, _ := c.AddFunc(schedule, func() {
				t.SendNotification(deliverMessage)
			})
			fmt.Println("id: ", id)
		}
		c.Start()
		t.reminderState.Store(methodName, c)
		fmt.Println("preStat", preStat)
	} else if setReminders.ReminderState == tp.OFF {
		preCron, _ := t.reminderState.Load(methodName)
		c := preCron.(*cron.Cron)
		fmt.Println("Entries: ", c.Entries())
		fmt.Println("TP OFF: Is going to STOP!! 2")
		fmt.Println("preStat", preStat)
		c.Remove(1)
		d := c.Stop()
		fmt.Println("STOP: ", d)
	} else {
		fmt.Println("TP HANG: Is going to HANG!! 3")
		fmt.Println(setReminders.ReminderState)
	}

	return nil
}

func (t telegram) SendNotification(msg tp.ReqMsg) error {
	var err error
	var response *http.Response
	emoji := ""
	if msg.State == pb.STATE_FAILURE {
		emoji = "❌"
	} else if msg.State == pb.STATE_SUCCESS {
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
		//else{
		//	log.Info().Str("module", "dispatcher").Msgf("dispatcher telegram response: %s", string(body))
		//}
	}
	return nil
}

func getUrl(token string) string {
	return fmt.Sprintf("https://api.telegram.org/bot%s", token)
}
