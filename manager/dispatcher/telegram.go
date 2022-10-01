package dispatcher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	tp "github.com/dsrvlabs/vatz/manager/types"
	"github.com/rs/zerolog/log"
)

var (
	Token  string
	ChatId string
)

// telegram: This is a sample code
// that helps to multi methods for notification.
type telegram struct {
	channel tp.Channel
	secret  string
	chatID  string
}

func (t telegram) SendTelegramNotification(text string) error {

	var err error
	var response *http.Response

	Token = t.secret
	ChatId = t.chatID

	url := fmt.Sprintf("%s/sendMessage", getUrl())
	body, _ := json.Marshal(map[string]string{
		"chat_id": ChatId,
		"text":    text,
	})

	response, err = http.Post(
		url,
		"application/json",
		bytes.NewBuffer(body),
	)

	if err != nil {
		log.Error().Str("module", "dispatcher").Msgf("dispatcher telegram Error: %s", err)
		return err
	}
	defer response.Body.Close()

	body, err = ioutil.ReadAll(response.Body)

	if err != nil {
		log.Error().Str("module", "dispatcher").Msgf("dispatcher telegram body parsing Error: %s", err)
		return err
	}
	// Log
	log.Info().Str("module", "dispatcher").Msgf("dispatcher telegram notification `%s` has sent ", text)
	log.Info().Str("module", "dispatcher").Msgf("dispatcher telegram response: %s", string(body))
	return nil
}

func (t telegram) SendNotification(request tp.ReqMsg) error {
	err := t.SendTelegramNotification(request.Msg)
	if err != nil {
		log.Error().Str("module", "dispatcher").Msgf("Sending a alert through Telegram has failed due to %s", err)
	}
	return nil
}

func getUrl() string {
	return fmt.Sprintf("https://api.telegram.org/bot%s", Token)
}
