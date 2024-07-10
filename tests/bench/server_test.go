package main

import (
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/marco-souza/marco.fly.dev/internal/server"
)

func BenchmarkProdServer(b *testing.B) {
	baseUrl := "https://marco.fly.dev"

	b.ResetTimer() // Reset the benchmark timer

	b.Run("sync GET /", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			http.Get(baseUrl + "/")
		}
	})

	b.Run("async GET /", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				http.Get(baseUrl + "/")
			}
		})
	})
}

func BenchmarkLocalServer(b *testing.B) {
	baseUrl := "http://localhost:3001"

	// set view folder
	b.Setenv("VIEWS", "../../views/")

	// disabling logs
	log.SetFlags(0)
	log.SetOutput(io.Discard)

	done := make(chan bool)
	s := server.New()

	go s.Start(&done)
	<-done

	b.ResetTimer() // Reset the benchmark timer

	b.Run("sync GET /", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			http.Get(baseUrl + "/")
		}
	})

	b.Run("async GET /", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				http.Get(baseUrl + "/")
			}
		})
	})
}
