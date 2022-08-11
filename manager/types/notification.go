package types

import (
	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"time"
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

type Field struct {
	Name   string `json:"name,omitempty"`
	Value  string `json:"value,omitempty"`
	Inline bool   `json:"inline,omitempty"`
}

// NotifyInfo contains detail dispatcher configs.
type NotifyInfo struct {
	Plugin     string            `json:"plugin"`
	Method     string            `json:"method"`
	Severity   pluginpb.SEVERITY `json:"severity"`
	State      pluginpb.STATE    `json:"state"`
	ExecuteMsg string            `json:"execute_msg"`
}
