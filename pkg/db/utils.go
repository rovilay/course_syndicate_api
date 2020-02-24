package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

// PopulateIndex ...
func PopulateIndex(c *mongo.Collection, key string, value int32, unique bool) (string, error) {
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)
	index := yieldIndexModel(key, value, unique)
	return c.Indexes().CreateOne(context.Background(), index, opts)
}

func yieldIndexModel(key string, value int32, unique bool) mongo.IndexModel {
	keys := bsonx.Doc{{Key: key, Value: bsonx.Int32(int32(value))}}
	index := mongo.IndexModel{}
	index.Keys = keys
	if unique {
		index.Options = options.Index().SetUnique(true)
	}
	return index
}
