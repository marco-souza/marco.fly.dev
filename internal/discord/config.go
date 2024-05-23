package discord

import "github.com/marco-souza/marco.fly.dev/internal/env"

type discordConfig struct {
	// app
	AppID     string
	AuthURL   string
	PublicKey string
	// bot
	ClientID     string
	ClientSecret string
	BotToken     string
}

func discordConfigLoad() *discordConfig {
	return &discordConfig{
		AuthURL:      env.Env("DISCORD_AUTH_URL", ""),
		AppID:        env.Env("DISCORD_APP_ID", ""),
		PublicKey:    env.Env("DISCORD_PUBLIC_KEY", ""),
		ClientID:     env.Env("DISCORD_CLIENT_ID", ""),
		ClientSecret: env.Env("DISCORD_CLIENT_SECRET", ""),
		BotToken:     env.Env("DISCORD_BOT_TOKEN", ""),
	}
}

var cfg = discordConfigLoad()
