package discord

import (
	"github.com/Shopify/go-lua"
)

func sendMsgWrapper(s *lua.State) int {
	// get channel from lua
	channel, ok := s.ToString(1) // {channel, message}
	if !ok {
		logger.Println("failed to get channel", channel)
		return 0
	}

	// get message from lua
	message, ok := s.ToString(2) // {channel, message}
	if !ok {
		logger.Println("failed to get message", message)
		s.PushBoolean(false) // {false, channel, message}
		return 1             // number of results
	}

	logger.Printf("sending message to channel: (%s) %s", channel, message)
	if err := DiscordService.SendMessage(channel, message); err != nil {
		logger.Fatalf("failed to send message: (%s) %s", channel, message)
		s.PushBoolean(false) // {false, channel, message}
		return 1             // number of results
	}

	logger.Println("message sent!")
	s.PushBoolean(true) // {true, channel, message}
	return 1            // number of results
}

func (d *discordService) PushClientLuaStack(l *lua.State) error {
	// ref: https://stackoverflow.com/a/37874926
	l.NewTable() // {}

	l.PushString("send_message")     // {}, "send_message"
	l.PushGoFunction(sendMsgWrapper) // {}, "send_message", sendMsgWrapper
	l.SetTable(-3)                   // {send_message: sendMsgWrapper}

	l.PushString("auth_url")  // {send_message: sendMsgWrapper}, "auth_url"
	l.PushString(cfg.AuthURL) // {send_message: sendMsgWrapper}, "auth_url", cfg.auth_url
	l.SetTable(-3)            // {send_message: sendMsgWrapper, auth_url: cfg.auth_url}

	// make it available globaly
	l.SetGlobal("discord")
	return nil
}
