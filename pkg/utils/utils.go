package utils

import (
	"encoding/json"
	"log"
	"net/http"
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

// ErrorHandler handles https error responses
func ErrorHandler(err *ErrorWithStatusCode, res http.ResponseWriter) {
	// log.Fatal(err.ErrorMessage.Error())
	res.Header().Set("Content-Type", "application/json")

	e, _ := json.Marshal(ErrorResponse{err.ErrorMessage.Error()})

	res.WriteHeader(err.StatusCode)
	res.Write(e)
}

// JSONResponseHandler handles http response in json
func JSONResponseHandler(res http.ResponseWriter, statusCode int, payload interface{}) {
	response, _ := json.Marshal(payload)

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(statusCode)
	res.Write(response)
}
