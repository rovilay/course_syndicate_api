package utils

import (
	"github.com/dgrijalva/jwt-go"
)

// ErrorWithStatusCode : This is error model.
type ErrorWithStatusCode struct {
	StatusCode   int
	ErrorMessage error
	Errors       interface{}
}

// ErrorResponse ...
type ErrorResponse struct {
	Message string      `json:"message"`
	Errors  interface{} `json:"errors,omitempty"`
}

// TokenPayload ...
type TokenPayload struct {
	ID    string `json:"id" bson:"id"`
	Email string `json:"email" bson:"email"`
}

// JWTClaims ...
type JWTClaims struct {
	ID    string `json:"id" bson:"id"`
	Email string `json:"email" bson:"email"`
	jwt.StandardClaims
}

// ContextKey ...
type ContextKey string

// SMTPServer ...
type SMTPServer struct {
	Host string
	Port string
}

// MailTemplateData ...
type MailTemplateData struct {
	Username string
	// Title    string
	// CourseID    string
	CourseTitle string
	// ModuleID    string
	ModuleTitle string
	ModuleLink  string
}
