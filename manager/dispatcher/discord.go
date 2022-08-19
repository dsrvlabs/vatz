package dispatcher

import (
	"bytes"
	"encoding/json"
	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/dsrvlabs/vatz/manager/config"
	tp "github.com/dsrvlabs/vatz/manager/types"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
	"time"
)

// DiscordColor describes color codes which are using for Discord msg.

type DiscordColor int

const (
	discordRed    tp.DiscordColor = 15548997
	discordYellow tp.DiscordColor = 16705372
	discordGreen  tp.DiscordColor = 65340
	discordGray   tp.DiscordColor = 9807270
	discordBlue   tp.DiscordColor = 4037805

	discordWebhookFormat string = "https://discord.com/api/webhooks/"
)

type discord struct{}

func (d discord) SendNotification(request tp.ReqMsg) error {
	cfg := config.GetConfig()

	err := d.sendNotificationForDiscord(request, cfg.Vatz.NotificationInfo.DiscordSecret)
	if err != nil {
		panic(err)
	}
	return nil
}

func (d discord) sendNotificationForDiscord(msg tp.ReqMsg, webhook string) error {
	if msg.ResourceType == "" {
		msg.ResourceType = "No Resource Type"
	}
	if msg.Msg == "" {
		msg.Msg = "No Message"
	}

	// Check discord secret
	if strings.Contains(webhook, discordWebhookFormat) {
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
		sMsg.Embeds[0].Fields = []tp.Field{{Name: msg.ResourceType, Value: msg.Msg, Inline: false}}
		sMsg.Embeds[0].Timestamp = time.Now()

		message, err := json.Marshal(sMsg)
		if err != nil {
			return err
		}

		req, _ := http.NewRequest("POST", webhook, bytes.NewReader(message))
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
