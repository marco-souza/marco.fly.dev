package main

import (
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/marco-souza/marco.fly.dev/internal/server"
)

var baseUrl = "https://marco.fly.dev"

func BenchmarkServer(b *testing.B) {
	b.Setenv("VIEWS", "../../views/")

	// disabling logs
	log.SetFlags(0)
	log.SetOutput(io.Discard)

	done := make(chan bool)
	s := server.New()

	go s.Start(&done)
	<-done

	b.ResetTimer() // Reset the benchmark timer

	b.Run("prod sync GET /", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			http.Get(baseUrl + "/")
		}
	})

	b.Run("prod async GET /", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				http.Get(baseUrl + "/")
			}
		})
	})

	baseUrl = "http://localhost:3001"

	b.Run("local sync GET /", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			http.Get(baseUrl + "/")
		}
	})

	b.Run("local async GET /", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				http.Get(baseUrl + "/")
			}
		})
	})
}
