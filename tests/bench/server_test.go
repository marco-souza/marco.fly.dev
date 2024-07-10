package main

import (
	"net/http"
	"testing"

	"github.com/marco-souza/marco.fly.dev/internal/env"
)

var routes = []string{
	"/",
	"/resume",
	"/blog",
	"/login",
}

func BenchmarkServerSync(b *testing.B) {
	for _, route := range routes {
		b.Run(route, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				fetchPage(route)
			}
		})
	}
}

func BenchmarkServerParallel(b *testing.B) {
	for _, route := range routes {
		b.Run(route, func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					fetchPage(route)
				}
			})
		})
	}
}

func fetchPage(route string) {
	baseUrl := env.Env("URL", "https://marco.fly.dev")
	http.Get(baseUrl + route)
}
