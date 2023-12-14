package main

import (
	"strings"
	"testing"

	"github.com/marco-souza/marco.fly.dev/internal/github"
)

func TestReturnedUrl(t *testing.T) {
	auth := github.Auth{}

	t.Run("validate redirect url", func(t *testing.T) {
		origin := "origin-url"
		redirectUrl := auth.RedirectLink(origin)

		if !strings.Contains(redirectUrl, "api.github") {
			t.Fatal("does not contain base api", redirectUrl)
		}
		if !strings.Contains(redirectUrl, origin) {
			t.Fatal("does not contain origin", redirectUrl)
		}
		if !strings.Contains(redirectUrl, "/oauth/authorize") {
			t.Fatal("is not an authorization endpoint", redirectUrl)
		}
		if !strings.Contains(redirectUrl, "state=") {
			t.Fatal("does not contain state", redirectUrl)
		}
		if !strings.Contains(redirectUrl, "scope=") {
			t.Fatal("does not contain scope", redirectUrl)
		}
	})
}
