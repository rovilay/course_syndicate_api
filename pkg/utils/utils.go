package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
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

	e, _ := json.Marshal(ErrorResponse{
		Message: err.ErrorMessage.Error(),
		Errors:  err.Errors,
	})

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

// GenerateToken ...
func GenerateToken(tp *TokenPayload) (string, error) {
	secret := EnvOrDefaultString("SECRET", "")
	et := EnvOrDefaultString("TOKEN_EXPIRATION_IN_HOURS", "24")
	i, err := strconv.Atoi(et)

	if err != nil {
		i = 24
	}

	expirationTime := time.Now().Add(time.Duration(i) * time.Hour)
	claims := &JWTClaims{
		ID:    tp.ID,
		Email: tp.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString([]byte(secret))

	return tokenString, err
}

// DecodeToken ...
func DecodeToken(token string) (*JWTClaims, error) {
	secret := EnvOrDefaultString("SECRET", "")
	claims := &JWTClaims{}

	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return claims, err
	}

	if !tkn.Valid {
		return claims, errors.New("Invalid token")
	}

	return claims, err
}

// ReadFile ...
func ReadFile(filename string) ([]byte, error) {
	file, err := os.Open(filename)

	if err != nil {
		return nil, err
	}

	fmt.Println("Successfully Opened " + filename)
	defer file.Close()

	return ioutil.ReadAll(file)
}
