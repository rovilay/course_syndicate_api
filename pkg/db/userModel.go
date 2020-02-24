package db

import (
	"time"

	root "github.com/rovilay/course_syndicate_api/pkg"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// UserModel ...
type UserModel struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id, omitempty"`
	FirstName    string             `json:"firstname" bson:"firstname"`
	LastName     string             `json:"lastname" bson:"lastname"`
	Email        string             `json:"email" bson:"email"`
	PasswordHash string             `json:"password" bson:"password"`
}

// CreateUserModel ...
func CreateUserModel(u *root.User) (newUser *UserModel, err error) {
	newUser = &UserModel{
		ID:        primitive.NewObjectIDFromTimestamp(time.Now()),
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
	}
	err = newUser.SetHashPassword(u.Password)

	return newUser, err
}

// ComparePasswordHash ...
func (u *UserModel) ComparePasswordHash(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))

	return err == nil
}

// SetHashPassword ...
func (u *UserModel) SetHashPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	u.PasswordHash = string(hash[:])

	return nil
}
