package main

import (
	"github.com/octofoxio/sparkle"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	client, err := mongo.NewClient(options.Client().ApplyURI(sparkle.LocalMongoDBURL))
}
