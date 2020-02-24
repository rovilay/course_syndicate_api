package controllers

import "github.com/rovilay/course_syndicate_api/pkg/db"

// UserController ...
type UserController struct {
	userService *db.UserService
}

// UserResponse ...
type UserResponse struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
}
