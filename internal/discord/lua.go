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
		logger.Printf("failed to send message: (%s) %s", channel, message)
		s.PushBoolean(false) // {false, channel, message}
		return 1             // number of results
	}

	logger.Println("message sent!")
	s.PushBoolean(true) // {true, channel, message}
	return 1            // number of results
}

func isWorkDay(l *lua.State) int {
	weekend := []time.Weekday{time.Saturday, time.Sunday}
	weekDay := time.Now().Weekday()

	logger.Printf("checking if %s is a work day: %v", weekDay, weekend)

	for _, day := range weekend {
		if day == weekDay {
			l.PushBoolean(false)
			return 1
		}
	}

	l.PushBoolean(true)
	return 1
}

func (d *discordService) PushClient(l *lua.State) {
	// ref: https://stackoverflow.com/a/37874926
	l.NewTable() // {}

	l.PushString("send_message")     // {}, "send_message"
	l.PushGoFunction(sendMsgWrapper) // {}, "send_message", sendMsgWrapper
	l.SetTable(-3)                   // {send_message: sendMsgWrapper}

	l.PushString("is_work_day") // {send_message: sendMsgWrapper}, "is_work_day"
	l.PushGoFunction(isWorkDay) // {send_message: sendMsgWrapper}, "is_work_day", isWorkDay
	l.SetTable(-3)              // {send_message: sendMsgWrapper, is_work_day: isWorkDay}

	l.PushString("auth_url")
	l.PushString(cfg.AuthURL)
	l.SetTable(-3)

	// make it available globaly
	l.SetGlobal("discord")
}
