package middlewares

import "course_syndicate_api/pkg/db"

// Validator ...
type Validator struct {
	Message                   string
	Errors                    map[string]string
	userService               *db.Service
	courseService             *db.Service
	courseSubscriptionService *db.Service
}
