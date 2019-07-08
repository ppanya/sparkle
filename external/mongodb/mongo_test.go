package mongodb

import (
	"context"
	"github.com/octofoxio/sparkle"
	"github.com/stretchr/testify/assert"
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

		err = db.MongoDB.Drop(context.Background())
		assert.NoError(t, err)
		assert.EqualValues(t, bb.ID, "test")
	})

	t.Run("transactional testing", func(t *testing.T) {
		db := NewLocal("integration-test")
		col := db.MongoDB.Collection("tx-test")

		t.Run("must able to commit", func(t *testing.T) {
			_ = col.Drop(context.Background())
			_ = db.Save(context.Background(), "tx-test", "test-pre", &B{
				A: A{Name: "jackee"},
			})
			provider := NewMongoTransactionalProvider(db.MongoDB.Client())
			err := provider.Begin(context.Background(), func(context sparkle.TransactionalContext) error {
				err := db.Save(context, "tx-test", "test", &B{
					A: A{Name: "jack"},
				})
				assert.NoError(t, err)
				return context.Commit(context)
			})
			assert.NoError(t, err)

			var bb B
			err = db.FindByID(context.Background(), "tx-test", "test", &bb)
			assert.NoError(t, err)
			assert.EqualValues(t, "jack", bb.Name)

		})
		t.Run("must able to rollback", func(t *testing.T) {
			_ = col.Drop(context.Background())
			_ = db.Save(context.Background(), "tx-test", "test-pre", &B{
				A: A{Name: "jackee"},
			})
			provider := NewMongoTransactionalProvider(db.MongoDB.Client())
			err := provider.Begin(context.Background(), func(context sparkle.TransactionalContext) error {
				err := db.Save(context, "tx-test", "test", &B{
					A: A{Name: "jack"},
				})
				assert.NoError(t, err)
				return context.Rollback(context)
			})
			assert.NoError(t, err)

			var bb B
			err = db.FindByID(context.Background(), "tx-test", "test", &bb)
			assert.Error(t, err)

		})
	})
}
