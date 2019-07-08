package migrate

import (
	"context"
	"github.com/octofoxio/sparkle"
	"github.com/octofoxio/sparkle/external/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func createEnsureDocumentToCollection(cl *mongo.Collection) (err error) {
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

func MigrateMongoCollection(db *mongodb.MongoDatabase, config *sparkle.Config) (err error) {
	err = createEnsureDocumentToCollection(db.MongoDB.Collection(config.UserCollectionName))
	if err != nil {
		return err
	}
	err = createEnsureDocumentToCollection(db.MongoDB.Collection(config.IdentityCollectionName))
	if err != nil {
		return err
	}

	err = createEnsureDocumentToCollection(db.MongoDB.Collection(config.SessionCollectionName))
	if err != nil {
		return err
	}
	return nil
}
