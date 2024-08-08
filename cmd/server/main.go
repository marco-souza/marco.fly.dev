package main

import (
	"github.com/marco-souza/marco.fly.dev/internal/di"
	"github.com/marco-souza/marco.fly.dev/internal/server"
)

func main() {
	di.Injectable(server.New(nil))
	s := di.MustInject(server.Server{})

	if err := s.Run(); err != nil {
		di.Clean()
		panic(err)
	}
}
