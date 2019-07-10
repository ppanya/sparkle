package migrate

import (
	"context"
	"github.com/octofoxio/sparkle"
	"github.com/octofoxio/sparkle/external/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

func createEnsureDocumentToCollection(db *mongodb.MongoDatabase, collectionName string) migrateStep {
	return func() (err error) {
		cl := db.MongoDB.Collection(collectionName)
		_, err = cl.UpdateOne(context.Background(),
			map[string]string{
				"_id": ".collection_ensure_just_ignore_this",
			},
			bson.D{{
				"$set", map[string]string{
					"_id":         ".collection_ensure_just_ignore_this",
					"description": "ensure that this collection is exists, due to transactional operation require collection to exists",
				},
			}},
			options.Update().SetUpsert(true),
		)
		return err
	}
}

func createUniqueIndex(db *mongodb.MongoDatabase, collectionName string, fieldName string) migrateStep {
	return func() (err error) {
		cl := db.MongoDB.Collection(collectionName)
		_, err = cl.Indexes().CreateOne(context.Background(), mongo.IndexModel{
			Keys:    bsonx.Doc{{Key: fieldName, Value: bsonx.Int32(-1)}},
			Options: options.Index().SetUnique(true).SetName(fieldName).SetSparse(true),
		})
		return err
	}
}

func DropMongoCollection(db *mongodb.MongoDatabase, config *sparkle.Config) error {
	var (
		userCollection     = db.MongoDB.Collection(config.UserCollectionName)
		sessionCollection  = db.MongoDB.Collection(config.SessionCollectionName)
		identityCollection = db.MongoDB.Collection(config.IdentityCollectionName)
	)
	err := userCollection.Drop(context.Background())
	if err != nil {
		return err
	}
	err = sessionCollection.Drop(context.Background())
	if err != nil {
		return err
	}
	err = identityCollection.Drop(context.Background())
	if err != nil {
		return err
	}
	return nil
}

type migrateStep func() error

func runner(ms ...migrateStep) error {
	for _, m := range ms {
		if err := m(); err != nil {
			return err
		}
	}
	return nil
}

func MustMigrateMongoCollection(db *mongodb.MongoDatabase, config *sparkle.Config) {
	err := MigrateMongoCollection(db, config)
	if err != nil {
		panic(err)
	}
}
func MigrateMongoCollection(db *mongodb.MongoDatabase, config *sparkle.Config) (err error) {
	return runner(
		createEnsureDocumentToCollection(db, config.UserCollectionName),
		createUniqueIndex(db, config.UserCollectionName, "lineid"),
		createUniqueIndex(db, config.UserCollectionName, "facebookid"),
		createUniqueIndex(db, config.UserCollectionName, "email"),

		createEnsureDocumentToCollection(db, (config.IdentityCollectionName)),

		createEnsureDocumentToCollection(db, (config.SessionCollectionName)),
	)
}
