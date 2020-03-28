package controllers

import "course_syndicate_api/pkg/db"

// UserController ...
type UserController struct {
	userService *db.Service
}

// CourseController ...
type CourseController struct {
	courseService       *db.Service
	courseModuleService *db.Service
}

// authResponse ...
type authResponse struct {
	Token string `json:"token,omitempty"`
}

type genericResponse struct {
	Message string `json:"message"`
}

type genericResponseWithData struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
