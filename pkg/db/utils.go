package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// IndexModel ...
type IndexModel struct {
	Keys   primitive.M
	Unique bool
}

// PopulateIndex ...
func PopulateIndex(c *mongo.Collection, indexes []mongo.IndexModel) ([]string, error) {
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)

	return c.Indexes().CreateMany(context.Background(), indexes, opts)
}

// YieldIndexes ...
func YieldIndexes(indexModels []IndexModel) []mongo.IndexModel {
	var indexes []mongo.IndexModel

	for _, m := range indexModels {
		index := mongo.IndexModel{}
		index.Keys = m.Keys
		if m.Unique {
			index.Options = options.Index().SetUnique(true)
		}

		indexes = append(indexes, index)
	}

	return indexes
}
