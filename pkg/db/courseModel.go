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
