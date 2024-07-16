package telegram

import "github.com/marco-souza/marco.fly.dev/internal/env"

var telegramChatID string
var telegramBotToken string
var telegramApiDomain string

func loadEnvs() {
	telegramChatID = env.Env("TELEGRAM_CHAT_ID", "161456907")
	telegramBotToken = env.Env("TELEGRAM_BOT_TOKEN", "")
	telegramApiDomain = env.Env(
		"TELEGRAM_API_DOMAIN", "api.telegram.org",
	)
}
