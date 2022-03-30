package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type reqMsg struct {
	FuncName     string `json:"func_name"`
	State        string `json:"state"`
	Msg          string `json:"msg"`
	Severity     string `json:"severity"`
	ResourceType string `json:"resource_type"`
}

type embed struct {
	Author struct {
		Name    string `json:"name,omitempty"`
		URL     string `json:"url,omitempty"`
		IconURL string `json:"icon_url,omitempty"`
	} `json:"author,omitempty"`
	Title       string `json:"title"`
	URL         string `json:"url,omitempty"`
	Description string `json:"description"`
	Color       int    `json:"color"`
	Fields      []struct {
		Name   string `json:"name,omitempty"`
		Value  string `json:"value,omitempty"`
		Inline bool   `json:"inline,omitempty"`
	} `json:"fields,omitempty"`
	Thumbnail struct {
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

type alertStatus struct{}

const (
	Fail    alertStatus = "FAIL"
	Success alertStatus = "SUCCESS"
)

var (
	webhookURL string
)

func init() {
	webhookURL = "https://discord.com/api/webhooks/955326233057566750/e-Tkq6mvxXYsx3Rqi3ttxDes76UXPpgk76Bnz5jF8_DwWgdNF-iNf7ZWdqi1nQnHO-zC"
}

func discord(request string) {
	message := makeForm(request)
	sendMsg(message)
}

func makeForm(request string) string {
	rMsg := reqMsg{}
	err := json.Unmarshal([]byte(request), &rMsg)
	if err != nil {
		fmt.Println("ERROR | Format is incorrect")
		panic(err)
	}

	sMsg := discordMsg{Embeds: make([]embed, 1)}
	if rMsg.State == Fail {
		fmt.Println("INFO | Got a fail status")
		sMsg.Embeds[0].Title = rMsg.Severity
		sMsg.Embeds[0].Description = rMsg.Msg
		sMsg.Embeds[0].Color = 15258703
		fmt.Println("INFO | Create fail message")
	} else if rMsg.State == Success {
		fmt.Println("INFO | Got a success status")
	}

	message, _ := json.Marshal(sMsg)
	return string(message)
}

func sendMsg(message string) {
	req, _ := http.NewRequest("POST", webhookURL, bytes.NewBufferString(message))
	req.Header.Set("Content-Type", "application/json")
	c := &http.Client{}
	_, err := c.Do(req)
	if err != nil {
		fmt.Println("ERROR | Failed to send discord message")
		panic(err)
	}
	fmt.Println("INFO | Send a discord")
}
