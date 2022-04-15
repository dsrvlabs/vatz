package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	pluginpb "github.com/xellos00/dk-yuba-proto/dist/proto/vatz/plugin/v1"
	"net/http"
	message "vatz/manager/model"
)

type notification struct {
}

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
	Title       string  `json:"title"`
	URL         string  `json:"url,omitempty"`
	Description string  `json:"description"`
	Color       int     `json:"color"`
	Fields      []field `json:"fields,omitempty"`
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
	sMsg := discordMsg{Embeds: make([]embed, 1)}
	sMsg.Embeds[0].Title = string(msg.Severity)
	sMsg.Embeds[0].Color = 15258703
	sMsg.Embeds[0].Fields = []field{{msg.FuncName, msg.Msg, false}}
	message, _ := json.Marshal(sMsg)
	req, _ := http.NewRequest("POST", webhook, bytes.NewBufferString(string(message)))
	req.Header.Set("Content-Type", "application/json")
	c := &http.Client{}
	_, err := c.Do(req)
	if err != nil {
		fmt.Println("ERROR | Failed to send discord message")
		panic(err)
	}
	return nil
}

func NewDispatcher() Notification {
	return &notification{}
}
