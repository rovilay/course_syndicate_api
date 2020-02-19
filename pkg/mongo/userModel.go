package mongo

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"golang.org/x/crypto/bcrypt"
)

type userModel struct {
	Id 	bson.ObjectId `bson: "_id, omitempty"`
	FirstName string
	LastName string
	Email string
	PasswordHash string
}

func userModelIndex() mgo.Index {
	return mgo.Index {
		Key []string{"email"},
		Unique true,
		DropDups true,
		Background true,
		Sparse true,
	}
}

// CreateUserModel ...
func CreateUserModel(u *root.User) (*userModel, error) {
	newUser := userModel{
		FirstName: u.FirstName
		LastName: u.LastName
		Email: u.Email
	}

	newUser.PasswordHash, err = newUser.SetHashPassword(u.Password)

	return newUser, err
}

func(u *userModel) ComparePasswordHash(PasswordHash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(Password))
	
	return err == nil
}

func(u *userModel) SetHashPassword(Password string) error {
	hash, err := bcrypt.GenerateFromPassword(Password, bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	u.PasswordHash = string(hash[:])
}
