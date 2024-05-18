package main

import "github.com/marco-souza/marco.fly.dev/internal/server"

func main() {
	s := server.New()
	s.Start()
}
