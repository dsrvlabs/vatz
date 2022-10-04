package dispatcher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/robfig/cron/v3"
	"net/http"
	"strings"
	"sync"
	"time"

	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	tp "github.com/dsrvlabs/vatz/manager/types"
	"github.com/rs/zerolog/log"
)

type DiscordColor int

const (
	discordRed    tp.DiscordColor = 15548997
	discordYellow tp.DiscordColor = 16705372
	discordGreen  tp.DiscordColor = 65340
	discordGray   tp.DiscordColor = 9807270
	discordBlue   tp.DiscordColor = 4037805

	discordWebhookFormat string = "https://discord.com/api/webhooks/"
)

type discord struct {
	host             string
	channel          tp.Channel
	secret           string
	reminderSchedule []string
	reminderState    sync.Map
}

func (d discord) SetDispatcher(firstRunMsg bool, preStat tp.StateFlag, notifyInfo tp.NotifyInfo) error {
	reqToNotify, setReminders, deliverMessage := messageHandler(firstRunMsg, preStat, notifyInfo)
	methodName := notifyInfo.Method
	fmt.Println(`==== START === (SetDispatcher) with method:`, methodName)

	if reqToNotify {
		d.SendNotification(deliverMessage)
	}
	fmt.Println("Map Value", d.reminderState)

	if setReminders.ReminderState == tp.ON {
		fmt.Println("TP ON: Is going to Start!! 1")
		dCron := &tp.CronTabSt{Crontab: cron.New(cron.WithLocation(time.UTC)), EntityID: 0}
		if _, ok := d.reminderState.Load(methodName); ok {
			fmt.Println("There's previous Cron")
			preStore, _ := d.reminderState.Load(methodName)
			dCron = preStore.(*tp.CronTabSt)
		}
		for _, schedule := range d.reminderSchedule {
			id, _ := dCron.Crontab.AddFunc(schedule, func() {
				d.SendNotification(deliverMessage)
			})
			dCron.Update(int(id))
			fmt.Println("id: ", id)
		}
		dCron.Crontab.Start()
		d.reminderState.Store(methodName, dCron)
		fmt.Println("reminderState", d.reminderState)
		fmt.Println("preStat", preStat)
	} else if setReminders.ReminderState == tp.OFF {
		fmt.Println("reminderState", d.reminderState)
		preCron, _ := d.reminderState.Load(methodName)
		c := preCron.(*tp.CronTabSt)
		fmt.Println("Entries: ", c.Crontab.Entries())
		fmt.Println("TP OFF: Is going to STOP!! 2")
		fmt.Println("preStat", preStat)
		c.Crontab.Remove(cron.EntryID(c.EntityID))
		d := c.Crontab.Stop()
		fmt.Println("STOP: ", d)
	} else {

		fmt.Println("TP HANG: Is going to HANG!! 3")
		fmt.Println(setReminders.ReminderState)
	}

	fmt.Println(`==== END === (SetDispatcher) with method:`, methodName)
	fmt.Println("")
	return nil
}

func (d discord) SendNotification(msg tp.ReqMsg) error {
	if msg.ResourceType == "" {
		msg.ResourceType = "No Resource Type"
	}
	if msg.Msg == "" {
		msg.Msg = "No Message"
	}

	// Check discord secret
	if strings.Contains(d.secret, discordWebhookFormat) {
		sMsg := tp.DiscordMsg{Embeds: make([]tp.Embed, 1)}
		switch msg.Severity {
		case pluginpb.SEVERITY_CRITICAL:
			sMsg.Embeds[0].Color = discordRed
		case pluginpb.SEVERITY_WARNING:
			sMsg.Embeds[0].Color = discordYellow
		case pluginpb.SEVERITY_INFO:
			sMsg.Embeds[0].Color = discordBlue
		default:
			sMsg.Embeds[0].Color = discordGray
		}

		sMsg.Embeds[0].Title = msg.Severity.String()
		sMsg.Embeds[0].Fields = []tp.Field{{Name: "(" + d.host + ") " + msg.ResourceType, Value: msg.Msg, Inline: false}}
		sMsg.Embeds[0].Timestamp = time.Now()

		message, err := json.Marshal(sMsg)
		if err != nil {
			return err
		}

		req, _ := http.NewRequest("POST", d.secret, bytes.NewReader(message))
		req.Header.Set("Content-Type", "application/json")
		c := &http.Client{}
		_, err = c.Do(req)
		if err != nil {
			log.Error().Str("module", "dispatcher").Msgf("dispatcher ch:discord-Send notification error: %s", err)
		}
		// TODO: Should handle response status.
	} else {
		log.Error().Str("module", "dispatcher").Msg("dispatcher ch:discord-Invalid discord webhook address")
	}
	return nil

}
