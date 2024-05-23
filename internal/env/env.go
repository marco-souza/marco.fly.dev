package env

import "os"

func Env(varEnv string, defaultValue string) string {
	value := os.Getenv(varEnv)
	if value != "" {
		return value
	}

	return defaultValue
}
