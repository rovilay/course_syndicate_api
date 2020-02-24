package utils

// ErrorWithStatusCode : This is error model.
type ErrorWithStatusCode struct {
	StatusCode   int
	ErrorMessage error
}

// ErrorResponse
type ErrorResponse struct {
	Message string `json:"message"`
}
