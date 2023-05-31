package types

import (
	"time"

	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/robfig/cron/v3"
)

// DiscordColor is color for a discord alert.
type DiscordColor int

// ReqMsg is Setup message into GRPC Type.
type ReqMsg struct {
	FuncName     string                 `json:"func_name"`
	State        pluginpb.STATE         `json:"state"`
	Msg          string                 `json:"msg"`
	Severity     pluginpb.SEVERITY      `json:"severity"`
	ResourceType string                 `json:"resource_type"`
	Options      map[string]interface{} `json:"options"`
}

// UpdateState is to uptade the state of pluginpb.
func (r *ReqMsg) UpdateState(stat pluginpb.STATE) {
	r.State = stat
}

// UpdateSeverity is to uptade the severity of pluginpb.
func (r *ReqMsg) UpdateSeverity(sev pluginpb.SEVERITY) {
	r.Severity = sev
}

// UpdateMSG is to update message
func (r *ReqMsg) UpdateMSG(message string) {
	r.Msg = message
}

// DiscordMsg is type for sending messages to a discord.
type DiscordMsg struct {
	Username  string  `json:"username,omitempty"`
	AvatarURL string  `json:"avatar_url,omitempty"`
	Content   string  `json:"content,omitempty"`
	Embeds    []Embed `json:"embeds"`
}

// Embed is imformation for detail message.
type Embed struct {
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
	Fields      []Field      `json:"fields,omitempty"`
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

// StateFlag is type that indicates the status of the plugins.
type StateFlag struct {
	State    pluginpb.STATE    `json:"state"`
	Severity pluginpb.SEVERITY `json:"severity"`
}

// CronTabSt is crontab structure.
type CronTabSt struct {
	Crontab  *cron.Cron `json:"crontab"`
	EntityID int        `json:"entity_id"`
}

// Update is to update CronTabSt.
func (in *CronTabSt) Update(entity int) {
	in.EntityID = entity
}

// Field is a structure for embeds that can be omitted.
type Field struct {
	Name   string `json:"name,omitempty"`
	Value  string `json:"value,omitempty"`
	Inline bool   `json:"inline,omitempty"`
}

// NotifyInfo contains detail dispatcher configs.
type NotifyInfo struct {
	Plugin     string            `json:"plugin"`
	Method     string            `json:"method"`
	Address    string            `json:"address"`
	Port       int               `json:"port"`
	Severity   pluginpb.SEVERITY `json:"severity"`
	State      pluginpb.STATE    `json:"state"`
	ExecuteMsg string            `json:"execute_msg"`
}

// Channel types for dispatchers.
type Channel string

// the type of channel.
const (
	Discord   Channel = "DISCORD"
	Telegram  Channel = "TELEGRAM"
	PagerDuty Channel = "PAGERDUTY"
)

// Reminder is for reminnig alert
type Reminder string

// The type of Reminder.
const (
	ON   Reminder = "ON"
	HANG Reminder = "HANG"
	OFF  Reminder = "OFF"
)
