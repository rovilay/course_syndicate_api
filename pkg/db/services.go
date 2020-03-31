package db

import (
	"log"

	root "course_syndicate_api/pkg"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Service ...
type Service struct {
	Collection *mongo.Collection
}

// NewService ...
func NewService(c *mongo.Client, config *root.MongoConfig, collectionName string) *Service {
	// connect to DB and get course collection
	Collection := c.Database(config.DBName).Collection(collectionName)

	return &Service{Collection}
}

// NewUserService ...
func NewUserService(c *mongo.Client, config *root.MongoConfig) *Service {
	// connect to DB and get user collection
	Collection := c.Database(config.DBName).Collection("users")

	// index email to ensure it is unique across the collections
	keys := bson.M{"email": 1}
	emailIndex := IndexModel{keys, true}
	_, err := PopulateIndex(Collection, YieldIndexes([]IndexModel{emailIndex}))
	if err != nil {
		log.Println("[POPULATE USER INDEX] ", err)
	}

	return &Service{Collection}
}

// NewCourseSubService ...
func NewCourseSubService(c *mongo.Client, config *root.MongoConfig) *Service {
	// connect to DB and get course_subscriptions collection
	Collection := c.Database(config.DBName).Collection("course_subscriptions")

	// index user and course to ensure it is unique across the collections
	keys := bson.M{"userId": 1, "courseId": 1}
	userCourseIndex := IndexModel{keys, true}
	_, err := PopulateIndex(Collection, YieldIndexes([]IndexModel{userCourseIndex}))
	if err != nil {
		log.Println("[POPULATE SUB_SCHEDULE INDEX] ", err)
	}

	return &Service{Collection}
}
