package controllers

import "github.com/rovilay/course_syndicate_api/pkg/db"

// UserController ...
type UserController struct {
	userService *db.UserService
}

// UserResponse ...
type authResponse struct {
	Token string `json:"token,omitempty"`
}
