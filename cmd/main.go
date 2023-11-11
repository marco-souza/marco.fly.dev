package main

import "github.com/marco-souza/marco.fly.dev/cmd/server"

func main() {
	s := server.New()
	s.Start()
}
