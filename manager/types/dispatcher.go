package types

import (
	"time"

	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/robfig/cron/v3"
)

type DiscordColor int

//Let's Setup this message into GRPC Type
type ReqMsg struct {
	FuncName     string            `json:"func_name"`
	State        pluginpb.STATE    `json:"state"`
	Msg          string            `json:"msg"`
	Severity     pluginpb.SEVERITY `json:"severity"`
	ResourceType string            `json:"resource_type"`
}

func (r *ReqMsg) UpdateState(stat pluginpb.STATE) {
	r.State = stat
}

func (r *ReqMsg) UpdateSeverity(sev pluginpb.SEVERITY) {
	r.Severity = sev
}

func (r *ReqMsg) UpdateMSG(message string) {
	r.Msg = message
}

type DiscordMsg struct {
	Username  string  `json:"username,omitempty"`
	AvatarURL string  `json:"avatar_url,omitempty"`
	Content   string  `json:"content,omitempty"`
	Embeds    []Embed `json:"embeds"`
}

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

type StateFlag struct {
	State    pluginpb.STATE    `json:"state"`
	Severity pluginpb.SEVERITY `json:"severity"`
}

type CronTabSt struct {
	Crontab  *cron.Cron `json:"crontab"`
	EntityID int        `json:"entity_id"`
}

func (in *CronTabSt) Update(entity int) {
	in.EntityID = entity
}

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

// Channel types for dispatchers
type Channel string

const (
	Discord  Channel = "DISCORD"
	Telegram Channel = "TELEGRAM"
)

type Reminder string

const (
	ON   Reminder = "ON"
	HANG Reminder = "HANG"
	OFF  Reminder = "OFF"
)
