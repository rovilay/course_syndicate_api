package mongo

import "gopkg.in/mgo.v2"

// UserService ...
type UserService struct {
	collection *mgo.Collection
}

// NewUserService ...
func NewUserService(s *Session, config *root.MongoConfig) *UserService {
	// connect to DB and get user collection
	collection := s.DB(config.DBName).C("user")

	return &UserService{collection}
}

// CreateUser ...
func (u *UserService) CreateUser(u *root.User) error {
	userModel, err := CreateUserModel(u)

	if err !== nil {
		return err
	}

	return u.collection.Insert(userModel)
}
