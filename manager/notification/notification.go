package notification

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/dsrvlabs/vatz/manager/config"
)

/* TODO: Discussion.
We need to discuss about notificatino module.
As I see this code, notification itself is described is dispatcher
but dispatcher and notification module should be splitted into two part.
*/

// State defines VATZ state.
type State string

// States
const (
	None       = State("NONE")
	Pending    = State("PENDING")
	InProgress = State("INPROGRESS")
	Faliure    = State("FAIILURE")
	Timeout    = State("TIMEOUT")
	Success    = State("SUCCESS")
)

// Severity defines notification level.
type Severity string

// Severities
const (
	Unknown  = Severity("UNKNOWN")
	Warning  = Severity("WARNING")
	Error    = Severity("ERROR")
	Critical = Severity("CRITICAL")
	Info     = Severity("INFO")
	Ok       = Severity("OK")
)

// DiscordColor describes color codes which are using for Discord msg.
type DiscordColor int

const (
	discordRed    DiscordColor = 15548997
	discordYellow DiscordColor = 16705372
	discordGreen  DiscordColor = 65340
	discordGray   DiscordColor = 9807270
	discordBlue   DiscordColor = 4037805

	discordWebhookFormat string = "https://discord.com/api/webhooks/"
)

var (
	notifSingleton Notification
	notifOnce      sync.Once
)

// ReqMsg is request message to send notification.
type ReqMsg struct {
	FuncName     string   `json:"func_name"`
	State        State    `json:"state"`
	Msg          string   `json:"msg"`
	Severity     Severity `json:"severity"`
	ResourceType string   `json:"resource_type"`
}

type discordMsg struct {
	Username  string  `json:"username,omitempty"`
	AvatarURL string  `json:"avatar_url,omitempty"`
	Content   string  `json:"content,omitempty"`
	Embeds    []embed `json:"embeds"`
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

type field struct {
	Name   string `json:"name,omitempty"`
	Value  string `json:"value,omitempty"`
	Inline bool   `json:"inline,omitempty"`
}

// NotifyInfo contains detail notification configs.
type NotifyInfo struct {
	Plugin     string            `json:"plugin"`
	Method     string            `json:"method"`
	Severity   pluginpb.SEVERITY `json:"severity"`
	State      pluginpb.STATE    `json:"state"`
	ExecuteMsg string            `json:"execute_msg"`
}

// Notification provides interfaces to send alert notification message with variable channel.
type Notification interface {
	SendDiscord(msg ReqMsg, webhook string) error
	GetNotifyInfo(response *pluginpb.ExecuteResponse, pluginName string, methodName string) NotifyInfo
	SendNotification(request ReqMsg) error
}

type notification struct {
}

func (d notification) SendNotification(request ReqMsg) error {
	cfg := config.GetConfig()

	err := d.SendDiscord(request, cfg.Vatz.NotificationInfo.DiscordSecret)
	if err != nil {
		panic(err)
	}

	return nil
}

func (d notification) GetNotifyInfo(response *pluginpb.ExecuteResponse, pluginName string, methodName string) NotifyInfo {
	notifyInfo := NotifyInfo{
		Plugin:     pluginName,
		Method:     methodName,
		Severity:   response.GetSeverity(),
		State:      response.GetState(),
		ExecuteMsg: response.GetMessage(),
	}

	return notifyInfo
}

func (d notification) SendDiscord(msg ReqMsg, webhook string) error {
	// Check empty contents
	if msg.Severity == "" {
		msg.Severity = Unknown
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
		case Critical:
			sMsg.Embeds[0].Color = discordRed
		case Warning:
			sMsg.Embeds[0].Color = discordYellow
		case Ok:
			sMsg.Embeds[0].Color = discordGreen
		case Info:
			sMsg.Embeds[0].Color = discordBlue
		default:
			sMsg.Embeds[0].Color = discordGray
		}
		sMsg.Embeds[0].Title = string(msg.Severity)
		sMsg.Embeds[0].Fields = []field{{Name: msg.ResourceType, Value: msg.Msg, Inline: false}}
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
			log.Println("ERROR | Failed to send discord message")
		}

		// TODO: Should handle response status.
	} else {
		log.Println("ERROR | Invalid discord webhook address")
	}
	return nil
}

// GetDispatcher create new notification dispatcher.
func GetDispatcher() Notification {
	notifOnce.Do(func() {
		notifSingleton = &notification{}
	})

	return notifSingleton
}
