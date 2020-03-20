package db

import (
	"time"

	root "github.com/rovilay/course_syndicate_api/pkg"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CourseModel ...
type CourseModel struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id, omitempty"`
	FirstName    string             `json:"firstname" bson:"firstname"`
	LastName     string             `json:"lastname" bson:"lastname"`
	Email        string             `json:"email" bson:"email"`
	PasswordHash string             `json:"password" bson:"password"`
}

// CreateCourseModel ...
func CreateCourseModel(u *root.User) (newUser *UserModel, err error) {
	newUser = &UserModel{
		ID:        primitive.NewObjectIDFromTimestamp(time.Now()),
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
	}
	err = newUser.SetHashPassword(u.Password)

	return newUser, err
}
