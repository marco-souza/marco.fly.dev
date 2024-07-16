package telegram

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/Shopify/go-lua"
)

var logger = slog.With("service", "telegram")

func Start() {
	logger.Info("starting telegram service")
	loadEnvs()
}

func Stop() {
	logger.Info("stopping telegram service")
}

func SendChatMessage(message string) error {
	telegramUrl := url.URL{
		Scheme: "https",
		Host:   telegramApiDomain,
		Path:   fmt.Sprintf("bot%s/sendMessage", telegramBotToken),
	}

	// prepare query params
	params := telegramUrl.Query()

	params.Add("parse_mode", "Markdown")
	params.Add("chat_id", telegramChatID)
	params.Add("text", message)

	// add query params back to request
	telegramUrl.RawQuery = params.Encode()

	strUrl := telegramUrl.String()
	logger.Info("calling telegram api", "url", strUrl)

	res, err := http.Get(strUrl)
	if err != nil {
		logger.Error("error calling telegram api", "err", err)
		return err
	}

	if res.StatusCode != 200 {
		err = fmt.Errorf("error calling telegram api: status: %d  url: %s", res.StatusCode, strUrl)
		logger.Error("error calling telegram api", "err", err)
		return err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Error("error reading telegram response", "err", err)
		return err
	}

	logger.Info(string(body))

	return nil
}

func sendMsgWrapper(s *lua.State) int {
	// get channel from lua
	message, ok := s.ToString(1) // {channel, message}
	if !ok {
		logger.Error("failed to get message", "message", message)
		s.PushBoolean(false) // {false, channel, message}
		return 1
	}

	logger.Info("sending message", "message", message)

	if err := SendChatMessage(message); err != nil {
		logger.Info("failed to send message", "message", message)
		s.PushBoolean(false) // {false, channel, message}
		return 1             // number of results
	}

	logger.Info("message sent", "message", message)
	s.PushBoolean(true) // {true, channel, message}
	return 1            // number of results
}

func PushClient(l *lua.State) {
	l.NewTable() // {}

	l.PushString("send_message")     // {}, "send_message"
	l.PushGoFunction(sendMsgWrapper) // {}, "send_message", sendMsgWrapper
	l.SetTable(-3)                   // {send_message: sendMsgWrapper}

	// make it available globaly
	l.SetGlobal("telegram")
}
