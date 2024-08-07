package telegram

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/Shopify/go-lua"
	"github.com/marco-souza/marco.fly.dev/internal/di"
)

var logger = slog.With("service", "telegram")

type TelegramService struct {
	chatID string
	token  string
	domain string
}

func New() *TelegramService {
	loadEnvs()
	return &TelegramService{
		chatID: telegramChatID,
		token:  telegramBotToken,
		domain: telegramApiDomain,
	}
}

func (t *TelegramService) Start() error {
	return nil
}

func (t *TelegramService) Stop() error {
	return nil
}

func (t *TelegramService) SendChatMessage(message string) error {
	telegramUrl := url.URL{
		Scheme: "https",
		Host:   t.domain,
		Path:   fmt.Sprintf("bot%s/sendMessage", t.token),
	}

	// prepare query params
	params := telegramUrl.Query()

	params.Add("parse_mode", "Markdown")
	params.Add("chat_id", t.chatID)
	params.Add("text", message)

	// add query params back to request
	telegramUrl.RawQuery = params.Encode()

	strUrl := telegramUrl.String()
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

	logger.Info("telegram api response", "body", string(body))
	return nil
}

func PushClient(l *lua.State) {
	l.NewTable() // {}

	l.PushString("send_message")     // {}, "send_message"
	l.PushGoFunction(sendMsgWrapper) // {}, "send_message", sendMsgWrapper
	l.SetTable(-3)                   // {send_message: sendMsgWrapper}

	// make it available globaly
	l.SetGlobal("telegram")
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

	t, err := di.Inject(TelegramService{})
	if err != nil {
		logger.Error("failed to get telegram service", "err", err)
		s.PushBoolean(false) // {false, channel, message}
		return 1
	}

	if err := t.SendChatMessage(message); err != nil {
		logger.Info("failed to send message", "message", message)
		s.PushBoolean(false) // {false, channel, message}
		return 1             // number of results
	}

	logger.Info("message sent", "message", message)
	s.PushBoolean(true) // {true, channel, message}
	return 1            // number of results
}
