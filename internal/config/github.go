package config

import "github.com/marco-souza/marco.fly.dev/internal/env"

type Github struct {
	ClientId      string
	ClientSecret  string
	CallbackUrl   string
	RefreshUrl    string
	LogoutUrl     string
	SignInUrl     string
	DashboardPage string
	LoginPage     string
	Scope         string
}

func GithubLoad() Github {
	conf := Github{
		ClientId:      env.Env("GITHUB_CLIENT_ID", ""),
		ClientSecret:  env.Env("GITHUB_CLIENT_SECRET", ""),
		DashboardPage: "/app/",
		LoginPage:     "/login",
		CallbackUrl:   "/api/auth/github/callback",
		RefreshUrl:    "/api/auth/github/refresh",
		LogoutUrl:     "/api/auth/github/logout",
		SignInUrl:     "/api/auth/github",
		Scope:         "read:user",
	}

	// TODO: mock for testing
	// if conf.ClientId == "" || conf.ClientSecret == "" {
	// 	panic("github credentials not found")
	// }

	return conf
}
