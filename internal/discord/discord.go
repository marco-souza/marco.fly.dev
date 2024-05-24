package discord

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

var logger = log.New(os.Stdout, "discord: ", log.LstdFlags)

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
		logger.Printf("error sending message: %v", err)
		return err
	}

	logger.Printf("message sent: '%s'", message)
	return nil
}

func (d *discordService) Close() error {
	return d.session.Close()
}

func (d *discordService) Open() error {
	return d.session.Open()
}

var DiscordService = new()
