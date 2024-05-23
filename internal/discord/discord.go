package discord

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

type discordService struct {
	session *discordgo.Session
	logger  *log.Logger
}

func New() *discordService {
	logger := log.New(os.Stdout, "discord: ", log.LstdFlags)
	session, err := discordgo.New("Bot " + cfg.BotToken)
	if err != nil {
		panic(err)
	}

	return &discordService{session, logger}
}
