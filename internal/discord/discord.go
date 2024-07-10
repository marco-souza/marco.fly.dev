package discord

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

var logger = slog.With("discord")

type discordService struct {
	session *discordgo.Session
}

func new() *discordService {
	session, err := discordgo.New("Bot " + cfg.BotToken)
	if err != nil {
		panic(err)
	}

	return &discordService{session}
}

func (d *discordService) SendMessage(channel, message string) error {
	_, err := d.session.ChannelMessageSend(channel, message)
	if err != nil {
		logger.Error("error sending message", "err", err)
		return err
	}

	logger.Info("message sent", "message", message)
	return nil
}

func (d *discordService) Close() error {
	return d.session.Close()
}

func (d *discordService) Open() error {
	return d.session.Open()
}

var DiscordService = new()
