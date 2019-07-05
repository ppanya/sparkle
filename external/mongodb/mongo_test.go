package mongodb

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

type A struct {
	Name string `bson:",omitempty"`
}
type B struct {
	A  `bson:",inline,omitempty"`
	ID string `bson:"_id,omitempty"`
}

func TestMongoDatabase(t *testing.T) {

	if testing.Short() {
		t.Skip()
	}
	t.Run("perform basic ACID", func(t *testing.T) {
		db := NewLocal("me")

		err := db.Save(context.Background(), "strings", "test", &A{
			Name: "jack",
		})
		assert.NoError(t, err)

		var aa B
		err = db.FindByID(context.Background(), "strings", "test", &aa)
		assert.NoError(t, err)

		var bb B
		err = db.FindOne(context.Background(), "strings", &B{
			ID: "test",
		}, &bb)

		err = db.DB.Drop(context.Background())
		assert.NoError(t, err)
		assert.EqualValues(t, bb.ID, "test")
	})

	t.Run("transactional commit testing", func(t *testing.T) {
		db := NewLocal("integration-test")
		col := db.DB.Collection("tx-testing")
		_, err := col.InsertOne(context.Background(), &A{
			Name: "jack",
		})
		err = db.DB.Client().UseSession(context.Background(), func(sessionContext mongo.SessionContext) error {
			err := sessionContext.StartTransaction()
			if err != nil {
				panic(err)
			}
			_, err = col.InsertOne(sessionContext, &A{
				Name: "jack",
			})

			assert.NoError(t, err)
			err = sessionContext.CommitTransaction(context.Background())
			sessionContext.EndSession(context.Background())
			if err != nil {
				panic(err)
			}
			return nil
		})
		assert.NoError(t, err)

		t.Run("must has exists after commit", func(t *testing.T) {
			var aa B
			err = db.FindByID(context.Background(), "tx-testing", "test", &aa)
			assert.NoError(t, err)
		})

		//err = db.DB.Drop(context.Background())
		//assert.NoError(t, err)
	})

	t.Run("transactional rollback testing", func(t *testing.T) {
		db := NewLocal("integration-test")
		c := context.Background()
		txProvider := NewMongoTransactionalProvider(db.DB.Client())

		txCtx, err := txProvider.Begin(c)
		assert.NoError(t, err)

		err = db.Save(txCtx, "transactional-testing", "test", &A{
			Name: "jack",
		})
		assert.NoError(t, err)

		err = txCtx.Rollback()
		assert.NoError(t, err)

		t.Run("must not found result after rollback", func(t *testing.T) {
			var aa B
			err = db.FindByID(context.Background(), "transactional-testing", "test", &aa)
			assert.Error(t, err)
		})

		err = db.DB.Drop(context.Background())
		assert.NoError(t, err)
	})

}
