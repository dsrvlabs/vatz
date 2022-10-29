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

	pb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	tp "github.com/dsrvlabs/vatz/manager/types"
	"github.com/rs/zerolog/log"
)

type DiscordColor int

const (
	discordRed           tp.DiscordColor = 15548997
	discordYellow        tp.DiscordColor = 16705372
	discordGreen         tp.DiscordColor = 65340
	discordGray          tp.DiscordColor = 9807270
	discordWebhookFormat string          = "https://discord.com/api/webhooks/"
)

type discord struct {
	host             string
	channel          tp.Channel
	secret           string
	reminderSchedule []string
	reminderCron     *cron.Cron
	entry            sync.Map
}

func (d *discord) SetDispatcher(firstRunMsg bool, preStat tp.StateFlag, notifyInfo tp.NotifyInfo) error {
	reqToNotify, reminderState, deliverMessage := messageHandler(firstRunMsg, preStat, notifyInfo)
	methodName := notifyInfo.Method

	if reqToNotify {
		d.SendNotification(deliverMessage)
	}

	if reminderState == tp.ON {
		newEntries := []cron.EntryID{}
		/*
			In case of reminder has to keep but stateFlag has changed,
			e.g.) CRITICAL -> WARNING
			e.g.) ERROR -> INFO -> ERROR
		*/
		if entries, ok := d.entry.Load(methodName); ok {
			for _, entry := range entries.([]cron.EntryID) {
				d.reminderCron.Remove(entry)
			}
			d.reminderCron.Stop()
		}
		for _, schedule := range d.reminderSchedule {
			id, _ := d.reminderCron.AddFunc(schedule, func() {
				d.SendNotification(deliverMessage)
			})
			newEntries = append(newEntries, id)
		}
		d.entry.Store(methodName, newEntries)
		d.reminderCron.Start()
	} else if reminderState == tp.OFF {
		entries, _ := d.entry.Load(methodName)
		for _, entity := range entries.([]cron.EntryID) {
			{
				d.reminderCron.Remove(entity)
			}
			d.reminderCron.Stop()
		}
	}

	return nil
}

func (d *discord) SendNotification(msg tp.ReqMsg) error {
	if msg.ResourceType == "" {
		msg.ResourceType = "No Resource Type"
	}
	if msg.Msg == "" {
		msg.Msg = "No Message"
	}

	// Check discord secret
	if strings.Contains(d.secret, discordWebhookFormat) {
		sMsg := tp.DiscordMsg{Embeds: make([]tp.Embed, 1)}
		sMsg.Embeds[0].Color = discordGray
		emoji := "🚨"

		if msg.State == pb.STATE_SUCCESS {
			switch {
			case msg.Severity == pb.SEVERITY_CRITICAL:
				sMsg.Embeds[0].Color = discordRed
				emoji = "‼️"
			case msg.Severity == pb.SEVERITY_WARNING:
				sMsg.Embeds[0].Color = discordYellow
				emoji = "❗"
			case msg.Severity == pb.SEVERITY_INFO:
				sMsg.Embeds[0].Color = discordGreen
				emoji = "✅"
			}
		}

		sMsg.Embeds[0].Title = fmt.Sprintf(`%s %s`, emoji, msg.Severity.String())
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
			log.Error().Str("module", "dispatcher").Msgf("Channel(Discord): Send notification error: %s", err)
		}
		// TODO: Should handle response status.
	} else {
		log.Error().Str("module", "dispatcher").Msg("Channel(Discord): Connection failed due to Invalid discord webhook address")
	}
	return nil
}
