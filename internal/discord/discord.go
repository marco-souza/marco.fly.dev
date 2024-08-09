package discord

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

var logger = slog.With("service", "discord")

type DiscordService struct {
	session *discordgo.Session
}

func New() *DiscordService {
	session, err := discordgo.New("Bot " + cfg.BotToken)
	if err != nil {
		panic(err)
	}

	return &DiscordService{session}
}

func (d *DiscordService) SendMessage(channel, message string) error {
	_, err := d.session.ChannelMessageSend(channel, message)
	if err != nil {
		logger.Error("error sending message", "err", err)
		return err
	}

	logger.Info("message sent", "message", message)
	return nil
}

func (d *DiscordService) Stop() error {
	return d.session.Close()
}

func (d *DiscordService) Start() error {
	return d.session.Open()
}
