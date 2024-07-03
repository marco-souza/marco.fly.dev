package config

import (
	"fmt"
	"strconv"

	"github.com/marco-souza/marco.fly.dev/internal/env"
)

type Config struct {
	Hostname    string
	Port        string
	Env         string // development | production
	SqliteUrl   string
	DatabaseUrl string
	RateLimit   int
	ResumeURL   string
	Github      Github
}

func Load() *Config {
	rateLimitStr := env.Env("RATE_LIMIT", "15")
	rateLimit, err := strconv.ParseUint(rateLimitStr, 10, 32)
	if err != nil {
		fmt.Println("cannot parse rate limit {}")
	}

	return &Config{
		Hostname:    env.Env("HOST", "localhost"),
		Port:        env.Env("PORT", "3001"),
		Env:         env.Env("ENV", "development"),
		DatabaseUrl: env.Env("DB_URL", "./test.db"), // TODO: deprecate
		SqliteUrl:   env.Env("DB_URL", "./test.db"),
		ResumeURL:   env.Env("RESUME_URL", "https://raw.githubusercontent.com/marco-souza/resume/main/RESUME.md"),
		RateLimit:   int(rateLimit),
		Github:      GithubLoad(),
	}
}
