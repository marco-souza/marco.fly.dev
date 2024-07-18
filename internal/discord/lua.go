package discord

import (
	"time"

	"github.com/Shopify/go-lua"
)

func sendMsgWrapper(s *lua.State) int {
	// get channel from lua
	channel, ok := s.ToString(1) // {channel, message}
	if !ok {
		logger.Error("failed to get channel", "channel", channel)
		s.PushBoolean(false) // {false, channel, message}
		return 1             // number of results
	}

	// get message from lua
	message, ok := s.ToString(2) // {channel, message}
	if !ok {
		logger.Error("failed to get message", "message", message)
		s.PushBoolean(false) // {false, channel, message}
		return 1             // number of results
	}

	logger.Info("sending message to channel", "channel", channel, "message", message)

	if err := DiscordService.SendMessage(channel, message); err != nil {
		logger.Info("failed to send message", "channel", channel, "message", message)
		s.PushBoolean(false) // {false, channel, message}
		return 1             // number of results
	}

	logger.Info("message sent", "channel", channel, "message", message)
	s.PushBoolean(true) // {true, channel, message}
	return 1            // number of results
}

func isWorkDay(l *lua.State) int {
	weekend := []time.Weekday{time.Saturday, time.Sunday}
	weekDay := time.Now().Weekday()

	logger.Info("checking if today is a work day", "weekDay", weekDay, "weekend", weekend)

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
