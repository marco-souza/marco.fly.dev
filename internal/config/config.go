package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Hostname    string
	Port        string
	Env         string // development | production
	DatabaseUrl string
	RateLimit   int
	Github      Github
}

func Load() *Config {
	// load .env file
	godotenv.Load()

	rateLimitStr := env("RATE_LIMIT", "15")
	rateLimit, err := strconv.ParseUint(rateLimitStr, 10, 32)
	if err != nil {
		fmt.Println("cannot parse rate limit {}")
	}

	return &Config{
		Hostname:    env("HOST", "localhost"),
		Port:        env("PORT", "3001"),
		Env:         env("ENV", "development"),
		DatabaseUrl: env("DB_URL", "./test.db"),
		RateLimit:   int(rateLimit),
		Github:      GithubLoad(),
	}
}

func env(varEnv string, defaultValue string) string {
	value := os.Getenv(varEnv)
	if value != "" {
		return value
	}
	return defaultValue
}
