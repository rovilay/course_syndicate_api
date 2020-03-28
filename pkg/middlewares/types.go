package middlewares

import "course_syndicate_api/pkg/db"

// Validator ...
type Validator struct {
	Message     string
	Errors      map[string]string
	UserService *db.Service
}
