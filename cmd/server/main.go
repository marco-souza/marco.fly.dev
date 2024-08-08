package main

import (
	"github.com/marco-souza/marco.fly.dev/internal/di"
	"github.com/marco-souza/marco.fly.dev/internal/server"
)

func main() {
	di.Injectable(server.New(nil))

	if err := di.Run(server.Server{}); err != nil {
		di.Clean()
		panic(err)
	}
}
