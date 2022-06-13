package notification

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	message "github.com/dsrvlabs/vatz/manager/model"
)

type notification struct {
}

type DiscordColor int

const (
	discordRed    DiscordColor = 15548997
	discordYellow DiscordColor = 16705372
	discordGreen  DiscordColor = 65340
	discordGray   DiscordColor = 9807270
	discordBlue   DiscordColor = 4037805

	discordWebhookFormat string = "https://discord.com/api/webhooks/"
)

type field struct {
	Name   string `json:"name,omitempty"`
	Value  string `json:"value,omitempty"`
	Inline bool   `json:"inline,omitempty"`
}

type embed struct {
	Author struct {
		Name    string `json:"name,omitempty"`
		URL     string `json:"url,omitempty"`
		IconURL string `json:"icon_url,omitempty"`
	} `json:"author,omitempty"`
	Title       string       `json:"title"`
	URL         string       `json:"url,omitempty"`
	Timestamp   time.Time    `json:"timestamp"`
	Description string       `json:"description"`
	Color       DiscordColor `json:"color"`
	Fields      []field      `json:"fields,omitempty"`
	Thumbnail   struct {
		URL string `json:"url,omitempty"`
	} `json:"thumbnail,omitempty"`
	Image struct {
		URL string `json:"url,omitempty"`
	} `json:"image,omitempty"`
	Footer struct {
		Text    string `json:"text,omitempty"`
		IconURL string `json:"icon_url,omitempty"`
	} `json:"footer,omitempty"`
}

type discordMsg struct {
	Username  string  `json:"username,omitempty"`
	AvatarURL string  `json:"avatar_url,omitempty"`
	Content   string  `json:"content,omitempty"`
	Embeds    []embed `json:"embeds"`
}

type Notification interface {
	SendDiscord(msg message.ReqMsg, webhook string) error
	GetNotifyInfo(response *pluginpb.ExecuteResponse, pluginName string, methodName string) map[interface{}]interface{}
}

func (d notification) GetNotifyInfo(response *pluginpb.ExecuteResponse, pluginName string, methodName string) map[interface{}]interface{} {
	notifyInfo := make(map[interface{}]interface{})
	notifyInfo["severity"] = response.GetSeverity().String()
	notifyInfo["state"] = response.GetState().String()
	notifyInfo["method_name"] = methodName
	notifyInfo["execute_message"] = response.GetMessage()
	notifyInfo["plugin_name"] = pluginName

	return notifyInfo
}

func (d notification) SendDiscord(msg message.ReqMsg, webhook string) error {
	// Check empty contents
	if msg.Severity == "" {
		msg.Severity = "No Severity"
	}
	if msg.ResourceType == "" {
		msg.ResourceType = "No Resource Type"
	}
	if msg.Msg == "" {
		msg.Msg = "No Message"
	}

	// Check discord secret
	if strings.Contains(webhook, discordWebhookFormat) {
		sMsg := discordMsg{Embeds: make([]embed, 1)}
		switch msg.Severity {
		case message.Critical:
			sMsg.Embeds[0].Color = discordRed
		case message.Warning:
			sMsg.Embeds[0].Color = discordYellow
		case message.Ok:
			sMsg.Embeds[0].Color = discordGreen
		case message.Info:
			sMsg.Embeds[0].Color = discordBlue
		default:
			sMsg.Embeds[0].Color = discordGray
		}
		sMsg.Embeds[0].Title = string(msg.Severity)
		sMsg.Embeds[0].Fields = []field{{Name: msg.ResourceType, Value: msg.Msg, Inline: false}}
		sMsg.Embeds[0].Timestamp = time.Now()
		message, _ := json.Marshal(sMsg)
		req, _ := http.NewRequest("POST", webhook, bytes.NewBufferString(string(message)))
		req.Header.Set("Content-Type", "application/json")
		c := &http.Client{}
		_, err := c.Do(req)
		if err != nil {
			log.Println("ERROR | Failed to send discord message")
		}
	} else {
		log.Println("ERROR | Invalid discord webhook address")
	}
	return nil
}

func NewDispatcher() Notification {
	return &notification{}
}
