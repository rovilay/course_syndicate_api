package db

import (
	"log"

	root "course_syndicate_api/pkg"

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
	_, err := PopulateIndex(Collection, "email", 1, true)
	if err != nil {
		log.Println("[POPULATE USER INDEX] ", err)
	}

	return &Service{Collection}
}
