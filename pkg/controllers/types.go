package controllers

import (
	"time"

	"course_syndicate_api/pkg/db"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserController ...
type UserController struct {
	userService *db.Service
}

// CourseController ...
type CourseController struct {
	courseService       *db.Service
	courseModuleService *db.Service
}

// CourseWithModule ...
type courseWithModule struct {
	ID              primitive.ObjectID      `json:"_id,omitempty" bson:"_id, omitempty"`
	Title           string                  `json:"title" bson:"title"`
	NumberOfModules int                     `json:"numberOfModules" bson:"numberOfModules"`
	Modules         []*db.CourseModuleModel `json:"modules"`
	CreatedAt       time.Time               `json:"createdAt" bson:"createdAt"`
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
