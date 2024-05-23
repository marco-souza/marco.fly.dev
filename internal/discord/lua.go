package discord

import (
	"time"

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

func isWorkDay() bool {
	weekend := []time.Weekday{time.Saturday, time.Sunday}
	weekDay := time.Now().Weekday()

	logger.Printf("checking if %s is a work day: %v", weekDay, weekend)

	for _, day := range weekend {
		if day == weekDay {
			return false
		}
	}

	return true
}

func (d *discordService) PushClientLuaStack(l *lua.State) error {
	// ref: https://stackoverflow.com/a/37874926
	l.NewTable() // {}

	l.PushString("send_message")     // {}, "send_message"
	l.PushGoFunction(sendMsgWrapper) // {}, "send_message", sendMsgWrapper
	l.SetTable(-3)                   // {send_message: sendMsgWrapper}

	l.PushString("is_work_day") // {send_message: sendMsgWrapper}, "is_work_day"
	l.PushBoolean(isWorkDay())  // {send_message: sendMsgWrapper}, "is_work_day", bool
	l.SetTable(-3)              // {send_message: sendMsgWrapper, is_work_day: bool}

	l.PushString("auth_url")  // {send_message: sendMsgWrapper, is_work_day: bool}, "auth_url"
	l.PushString(cfg.AuthURL) // {send_message: sendMsgWrapper, is_work_day: bool}, "auth_url", cfg.auth_url
	l.SetTable(-3)            // {send_message: sendMsgWrapper, is_work_day: bool, auth_url: cfg.auth_url}

	// make it available globaly
	l.SetGlobal("discord")
	return nil
}
