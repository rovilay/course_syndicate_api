package controllers

import (
	"time"

	"course_syndicate_api/pkg/db"
	"course_syndicate_api/pkg/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserController ...
type UserController struct {
	userService *db.Service
}

// CourseController ...
type CourseController struct {
	courseService               *db.Service
	courseModuleService         *db.Service
	courseSubscriptionService   *db.Service
	subscriptionScheduleService *db.Service
	smtpService                 *utils.SMTPServer
}

type courseWithModule struct {
	ID              primitive.ObjectID      `json:"_id,omitempty" bson:"_id, omitempty"`
	Title           string                  `json:"title" bson:"title"`
	NumberOfModules int                     `json:"numberOfModules" bson:"numberOfModules"`
	Modules         []*db.CourseModuleModel `json:"modules" bson:"modules"`
	CreatedAt       time.Time               `json:"createdAt" bson:"createdAt"`
}

type courseSubscription struct {
	ID               string    `json:"_id,omitempty" bson:"_id, omitempty"`
	UserID           string    `json:"userId" bson:"userId"`
	CourseID         string    `json:"courseId" bson:"courseId"`
	ModulesCompleted int       `json:"modulesCompleted" bson:"modulesCompleted"`
	Schedule         []int64   `json:"schedule" bson:"schedule"`
	CreatedAt        time.Time `json:"createdAt" bson:"createdAt"`
}

type scheduleSubscription struct {
	ID               primitive.ObjectID `json:"_id,omitempty" bson:"_id, omitempty"`
	ModulesCompleted int                `json:"modulesCompleted" bson:"modulesCompleted"`
	CreatedAt        time.Time          `json:"createdAt" bson:"createdAt"`
}

type scheduleUser struct {
	FirstName string `json:"firstname" bson:"firstname"`
	LastName  string `json:"lastname" bson:"lastname"`
	Email     string `json:"email" bson:"email"`
}

type scheduleCourse struct {
	ID              primitive.ObjectID `json:"_id,omitempty" bson:"_id, omitempty"`
	Title           string             `json:"title" bson:"title"`
	NumberOfModules int                `json:"numberOfModules" bson:"numberOfModules"`
	CreatedAt       time.Time          `json:"createdAt" bson:"createdAt"`
}

type fetchSchedulesResult struct {
	ID            primitive.ObjectID      `json:"_id,omitempty" bson:"_id, omitempty"`
	Schedule      int64                   `json:"schedule" bson:"schedule"`
	Completed     bool                    `json:"completed" bson:"completed"`
	Subscription  scheduleSubscription    `json:"subscription" bson:"subscription"`
	User          scheduleUser            `json:"user" bson:"user"`
	Course        scheduleCourse          `json:"course" bson:"course"`
	CourseModules []*db.CourseModuleModel `json:"courseModules" bson:"courseModules"`
	CreatedAt     time.Time               `json:"createdAt" bson:"createdAt"`
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
