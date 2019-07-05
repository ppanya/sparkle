package mongodb

import (
	"context"
	"fmt"
	"github.com/octofoxio/sparkle"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoTransactionalContext struct {
	context.Context
	mongo.Session
}

func NewMongoTransactionalContext(context context.Context, session mongo.Session) *MongoTransactionalContext {
	return &MongoTransactionalContext{Context: context, Session: session}
}

func (m *MongoTransactionalContext) Rollback() error {
	err := m.AbortTransaction(m)
	if err != nil {
		return err
	}
	m.Session.EndSession(m)
	return nil
}

func (m *MongoTransactionalContext) Commit() error {
	err := m.CommitTransaction(m)
	if err != nil {
		return err
	}
	m.Session.EndSession(m)
	return nil
}

func NewMongoTransactionalProvider(c *mongo.Client) *MongoTransactionalProvider {
	return &MongoTransactionalProvider{c: c}
}

type MongoTransactionalProvider struct {
	c *mongo.Client
}

func (m *MongoTransactionalProvider) Begin(ctx context.Context) (sparkle.TransactionalContext, error) {
	s, err := m.c.StartSession()
	if err != nil {
		return nil, err
	}
	err = s.StartTransaction()
	if err != nil {
		return nil, err
	}
	return &MongoTransactionalContext{
		Context: ctx,
		Session: s,
	}, nil

}

type MongoCollection struct {
	CollectionName string
	DB             *MongoDatabase
}

func (m *MongoCollection) FindByID(ctx context.Context, ID string, value interface{}) error {
	return m.DB.FindByID(ctx, m.CollectionName, ID, value)
}

func (m *MongoCollection) Save(ctx context.Context, ID string, entity interface{}) error {
	return m.DB.Save(ctx, m.CollectionName, ID, entity)
}

func (m *MongoCollection) DeleteByID(ctx context.Context, ID string) error {
	panic("not implement")
}

type MongoDatabase struct {
	DB *mongo.Database
}

func (m *MongoDatabase) Collection(name string) sparkle.Collection {
	return &MongoCollection{
		DB:             m,
		CollectionName: name,
	}
}

func (m *MongoDatabase) FindOne(ctx context.Context, Collection string, filter interface{}, value interface{}) error {
	b, err := bson.Marshal(filter)
	if err != nil {
		return err
	}
	return m.DB.
		Collection(Collection).
		FindOne(
			ctx,
			b,
		).
		Decode(value)
}

func (m *MongoDatabase) FindByID(ctx context.Context, Collection, ID string, value interface{}) error {
	return m.DB.
		Collection(Collection).
		FindOne(ctx, bson.D{{
			"_id", ID,
		}}).
		Decode(value)
}

func (m *MongoDatabase) Save(ctx context.Context, Collection, ID string, entity interface{}) error {

	_, err := m.DB.
		Collection(Collection).
		UpdateOne(
			ctx,
			bson.D{{
				"_id", ID,
			}},
			bson.D{{
				"$set", entity,
			}},
			options.Update().SetUpsert(true),
		)
	return err
}

func (m *MongoDatabase) DeleteByID(ctx context.Context, Collection, ID string) error {
	panic("implement me")
}

func New(db *mongo.Database) *MongoDatabase {
	return &MongoDatabase{
		DB: db,
	}
}

func NewLocal(databaseName string) *MongoDatabase {
	client, err := mongo.NewClient(
		options.Client().ApplyURI(sparkle.LocalMongoDBURL))
	if err != nil {
		panic(err)
	}
	c, _ := context.WithTimeout(context.Background(), time.Minute*3)
	defer func() {
		c.Done()
		if c.Err() != nil {
			panic(fmt.Sprintf("maybe local mongodb is offline? (%s)", c.Err().Error()))
		}
	}()
	err = client.Connect(c)
	if err != nil {
		panic(err)
	}
	ss, err := client.ListDatabaseNames(context.Background(), bson.D{})
	if err != nil {
		panic(err)
	}
	fmt.Println(ss)
	return &MongoDatabase{
		DB: client.Database(databaseName),
	}
}
