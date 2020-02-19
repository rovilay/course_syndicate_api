package root

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// EnvOrDefaultString ...
func EnvOrDefaultString(envVar string, defaultValue string) string {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	value := os.Getenv(envVar)

	if value == "" {
		return defaultValue
	}

	return value
}
