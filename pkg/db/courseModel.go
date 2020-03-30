package db

import (
	"time"

	root "course_syndicate_api/pkg"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CourseModel ...
type CourseModel struct {
	ID              primitive.ObjectID `json:"_id,omitempty" bson:"_id, omitempty"`
	Title           string             `json:"title" bson:"title"`
	NumberOfModules int                `json:"numberOfModules" bson:"numberOfModules"`
	CreatedAt       time.Time          `json:"createdAt" bson:"createdAt"`
}

// CourseModuleModel ...
type CourseModuleModel struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id, omitempty"`
	Title     string             `json:"title" bson:"title"`
	CourseID  primitive.ObjectID `json:"courseId" bson:"courseId"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
}

// CourseSubscriptionModel ...
type CourseSubscriptionModel struct {
	ID               primitive.ObjectID `json:"_id,omitempty" bson:"_id, omitempty"`
	UserID           primitive.ObjectID `json:"userId" bson:"userId"`
	CourseID         primitive.ObjectID `json:"courseId" bson:"courseId"`
	ModulesCompleted int                `json:"modulesCompleted" bson:"modulesCompleted"`
	CreatedAt        time.Time          `json:"createdAt" bson:"createdAt"`
}

// SubscriptionScheduleModel ...
type SubscriptionScheduleModel struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id, omitempty"`
	SubscriptionID primitive.ObjectID `json:"subscriptionId" bson:"subscriptionId"`
	Schedule       int64              `json:"schedule" bson:"schedule"`
	Completed      bool               `json:"completed" bson:"completed"`
	CreatedAt      time.Time          `json:"createdAt" bson:"createdAt"`
}

// CourseWithModule ...
type CourseWithModule struct {
	ID              primitive.ObjectID  `json:"_id,omitempty" bson:"_id,omitempty"`
	Title           string              `json:"title" bson:"title"`
	NumberOfModules int                 `json:"numberOfModules" bson:"numberOfModules"`
	Modules         []CourseModuleModel `json:"modules,omitempty" bson:"modules,omitempty"`
	CreatedAt       time.Time           `json:"createdAt" bson:"createdAt"`
}

// CreateCourseModel ...
func CreateCourseModel(c *root.Course) *CourseModel {
	newCourse := &CourseModel{
		ID:              primitive.NewObjectIDFromTimestamp(time.Now()),
		Title:           c.Title,
		NumberOfModules: c.NumberOfModules,
		CreatedAt:       time.Now(),
	}

	return newCourse
}

// CreateCourseModuleModel ...
func CreateCourseModuleModel(cid primitive.ObjectID, cm *root.CourseModule) *CourseModuleModel {
	newCourseModule := &CourseModuleModel{
		ID:        primitive.NewObjectIDFromTimestamp(time.Now()),
		Title:     cm.Title,
		CourseID:  cid,
		CreatedAt: time.Now(),
	}

	return newCourseModule
}

// CreateCourseSubscriptionModel ...
func CreateCourseSubscriptionModel(uid, cid primitive.ObjectID, cs []int64) *CourseSubscriptionModel {
	newCourseSubscription := &CourseSubscriptionModel{
		ID:               primitive.NewObjectIDFromTimestamp(time.Now()),
		UserID:           uid,
		CourseID:         cid,
		ModulesCompleted: 0,
		CreatedAt:        time.Now(),
	}

	return newCourseSubscription
}

// CreateSubscriptionScheduleModel ...
func CreateSubscriptionScheduleModel(subscriptionID primitive.ObjectID, schedule int64) *SubscriptionScheduleModel {
	newSubscriptionSchedule := &SubscriptionScheduleModel{
		ID:             primitive.NewObjectIDFromTimestamp(time.Now()),
		SubscriptionID: subscriptionID,
		Schedule:       schedule,
		Completed:      false,
		CreatedAt:      time.Now(),
	}

	return newSubscriptionSchedule
}
