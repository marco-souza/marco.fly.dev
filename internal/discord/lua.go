package discord

import (
	"github.com/Shopify/go-lua"
)

func sendMsgWrapper(s *lua.State) int {
	// get channel from lua
	channel, ok := s.ToString(1)
	if !ok {
		return 0
	}

	// get message from lua
	message, ok := s.ToString(2)
	if !ok {
		return 0
	}

	logger.Printf("sending message to channel: (%s) %s", channel, message)
	if err := DiscordService.SendMessage(channel, message); err != nil {
		logger.Fatalf("failed to send message: (%s) %s", channel, message)
		s.PushBoolean(false)
		return 1 // number of results
	}

	logger.Println("message sent!")
	s.PushBoolean(true)
	return 1 // number of results
}

func (d *discordService) PushClientLuaStack(l *lua.State) error {
	// l.NewTable()
	//
	// 	l.PushString("discord")
	//
	// 		l.NewTable()
	// 		l.PushString("send_message")
	// 		l.PushGoFunction(sendMsgWrapper)
	//
	// 	l.SetTable(-3)
	//
	// l.SetTable(-3)
	// l.SetGlobal("api")
	l.NewTable()

	l.PushString("send_message")
	l.PushGoFunction(sendMsgWrapper)

	l.SetTable(-3)

	// make it available globaly
	l.SetGlobal("discord")
	return nil
}
