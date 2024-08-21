package dispatcher

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	tp "github.com/dsrvlabs/vatz/types"
	"github.com/dsrvlabs/vatz/utils"
	"net/http"
	"strings"
	"sync"
	"time"

	pb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
)

// DiscordColor is type for discord message color.
type DiscordColor int

const (
	discordRed    tp.DiscordColor = 15548997
	discordYellow tp.DiscordColor = 16705372
	discordGreen  tp.DiscordColor = 65340
	discordGray   tp.DiscordColor = 9807270
)

var discordWebhookFormats = []string{
	"https://discord.com/api/webhooks/",
	"https://discordapp.com/api/webhooks/",
}

type discord struct {
	host             string
	channel          tp.Channel
	secret           string
	notificationFlag string
	reminderSchedule []string
	reminderCron     *cron.Cron
	entry            sync.Map
}

func containsAny(s string, substrings []string) bool {
	for _, substr := range substrings {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}

func (d *discord) SetDispatcher(firstRunMsg bool, pluginNotificationFlag string, preStat tp.StateFlag, notifyInfo tp.NotifyInfo) error {
	reqToNotify, reminderState, deliverMessage := messageHandler(firstRunMsg, preStat, notifyInfo)
	pUnique := deliverMessage.Options["pUnique"].(string)
	flagEnabled, sameFlagExists := utils.IsNotifiedEnabledAndSend(d.notificationFlag, pluginNotificationFlag)
	if !flagEnabled || flagEnabled && sameFlagExists {
		if reqToNotify {
			err := d.SendNotification(deliverMessage)
			if err != nil {
				log.Error().Str("module", "dispatcher").Msgf("Channel(Discord): Send notification error: %s", err)
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
		if entries, ok := d.entry.Load(pUnique); ok {
			for _, entry := range entries.([]cron.EntryID) {
				d.reminderCron.Remove(entry)
			}
			d.reminderCron.Stop()
		}
		for _, schedule := range d.reminderSchedule {
			id, _ := d.reminderCron.AddFunc(schedule, func() {
				err := d.SendNotification(deliverMessage)
				if err != nil {
					log.Error().Str("module", "dispatcher").Msgf("Channel(Discord): Send notification error: %s", err)
				}
			})
			newEntries = append(newEntries, id)
		}
		d.entry.Store(pUnique, newEntries)
		d.reminderCron.Start()
	} else if reminderState == tp.OFF {
		entries, _ := d.entry.Load(pUnique)
		if _, ok := entries.([]cron.EntryID); ok {
			for _, entity := range entries.([]cron.EntryID) {
				d.reminderCron.Remove(entity)
			}
			d.reminderCron.Stop()
		}
	}

	return nil
}

func (d *discord) SendNotification(msg tp.ReqMsg) error {
	if containsAny(d.secret, discordWebhookFormats) {
		sendingMsg := tp.DiscordMsg{Embeds: make([]tp.Embed, 1)}
		sendingMsg.Embeds[0].Color = discordGray
		emoji := emojiER

		if msg.ResourceType == "" {
			msg.ResourceType = "No Resource Type"
		}

		if msg.Msg == "" {
			msg.Msg = "No Message"
		}

		if msg.State == pb.STATE_SUCCESS {
			switch {
			case msg.Severity == pb.SEVERITY_CRITICAL:
				sendingMsg.Embeds[0].Color = discordRed
				emoji = emojiDoubleEX
			case msg.Severity == pb.SEVERITY_WARNING:
				sendingMsg.Embeds[0].Color = discordYellow
				emoji = emojiSingleEx
			case msg.Severity == pb.SEVERITY_INFO:
				sendingMsg.Embeds[0].Color = discordGreen
				emoji = emojiCheck
			}
		}

		sendingMsg.Embeds[0].Title = fmt.Sprintf(`%s %s`, emoji, msg.Severity.String())
		sendingMsg.Embeds[0].Fields = []tp.Field{{Name: "(" + d.host + ") " + msg.ResourceType, Value: msg.Msg, Inline: false}}
		sendingMsg.Embeds[0].Timestamp = time.Now()

		message, err := json.Marshal(sendingMsg)
		if err != nil {
			return err
		}

		req, _ := http.NewRequest("POST", d.secret, bytes.NewReader(message))
		req.Header.Set("Content-Type", "application/json")
		c := &http.Client{}
		_, err = c.Do(req)
		if err != nil {
			log.Error().Str("module", "dispatcher").Msgf("Channel(Discord): Send notification error: %s", err)
			errorMessage := fmt.Sprintf("Channel(Discord) Error: %s", err)
			return errors.New(errorMessage)
		}
		// TODO: Should handle response status.
	} else {
		log.Error().Str("module", "dispatcher").Msg("Channel(Discord): Connection failed due to Invalid discord webhook address")
		return errors.New("Channel(Discord) Error: Invalid discord webhook address. ")
	}
	return nil
}
