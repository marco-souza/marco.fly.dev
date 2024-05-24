package env

import (
	"os"

	"github.com/joho/godotenv"
)

func Env(varEnv string, defaultValue string) string {
	// load .env file
	godotenv.Load()

	value := os.Getenv(varEnv)
	if value != "" {
		return value
	}

	return defaultValue
}
