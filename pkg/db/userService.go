package db

import (
	"log"

	root "github.com/rovilay/course_syndicate_api/pkg"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserService ...
type UserService struct {
	Collection *mongo.Collection
}

// NewUserService ...
func NewUserService(c *mongo.Client, config *root.MongoConfig) *UserService {
	// connect to DB and get user collection
	Collection := c.Database(config.DBName).Collection("users")
	_, err := PopulateIndex(Collection, "email", 1, true)
	if err != nil {
		log.Println("[POPULATE USER INDEX] ", err)
	}

	return &UserService{Collection}
}

// // NewUserService ...
// func NewUserService(s *mgo.Session, config *root.MongoConfig) *UserService {
// 	// connect to DB and get user collection
// 	collection := s.DB(config.DBName).C("user")

// 	return &UserService{collection}
// }

// CreateUser ...
// func (u *UserService) CreateUser(user *root.User) error {
// 	userModel, err := CreateUserModel(user)

// 	if err != nil {
// 		return err
// 	}

// 	result, err := u.collection.InsertOne(context.Background(), userModel)
// 	fmt.Println("[INSERT RESULT]", result)

// 	return err
// }
