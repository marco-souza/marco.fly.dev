package main

import (
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/marco-souza/marco.fly.dev/internal/server"
)

var routes = []string{
	"/",
	"/resume",
	"/blog",
	"/login",
}

func BenchmarkProdServerSync(b *testing.B) {
	for _, route := range routes {
		b.Run(route, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				fetchProdPage(route)
			}
		})
	}
}

func BenchmarkProdServerParallel(b *testing.B) {
	for _, route := range routes {
		b.Run(route, func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					fetchProdPage(route)
				}
			})
		})
	}
}

func BenchmarkLocalServerSync(b *testing.B) {
	teardown := localServerSetup(b)
	defer teardown()

	for _, route := range routes {
		b.Run(route, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				fetchLocalPage(route)
			}
		})
	}
}

func BenchmarkLocalServerParallel(b *testing.B) {
	teardown := localServerSetup(b)
	defer teardown()

	for _, route := range routes {
		b.Run(route, func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					fetchLocalPage(route)
				}
			})
		})
	}
}

func fetchProdPage(route string) {
	baseUrl := "https://marco.fly.dev"
	http.Get(baseUrl + route)
}

func fetchLocalPage(route string) {
	baseUrl := "http://localhost:3001"
	http.Get(baseUrl + route)
}

func localServerSetup(b *testing.B) func() {
	// set view folder
	b.Setenv("VIEWS", "../../views/")

	// disabling logs
	l := 1
	if &l == nil {
		log.SetFlags(0)
		log.SetOutput(io.Discard)
	}

	done := make(chan bool)
	s := server.New()

	go s.Start(&done)
	<-done

	fetchLocalPage("/") // warm up

	// Reset the benchmark timer
	b.ResetTimer()

	return s.Shutdown
}
