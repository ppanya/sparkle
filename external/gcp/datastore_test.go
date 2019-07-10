package gcp

import (
	"context"
	"github.com/octofoxio/sparkle"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func executeTime(name string) func() {
	start := time.Now()
	return func() {
		elapsed := time.Since(start)
		log.Printf("Datastore %s took %s", name, elapsed)
	}
}
func TestNewGCPDatastore(t *testing.T) {

	var db sparkle.Database = NewGCPDatastore("catcat-development", "spike-local")
	type Entity struct {
		Hi string
	}

	var done = executeTime("Save()")
	err := db.Save(context.Background(), "ping", "key-pong", &Entity{
		Hi: "Pong",
	})
	assert.NoError(t, err)
	done()

	var result Entity
	done = executeTime("FindByID()")
	err = db.FindByID(context.Background(), "ping", "key-pong", &result)
	assert.NotNil(t, result)
	t.Log(result.Hi)
	assert.NoError(t, err)
	done()

	done = executeTime("DeleteByID()")
	err = db.DeleteByID(context.Background(), "ping", "key-pong")
	assert.NoError(t, err)
	done()

	col := db.Collection("test-collection")
	err = col.Save(context.Background(), "test-item", &Entity{
		Hi: "hello",
	})
	assert.NoError(t, err)

	var bbb Entity

	//WIP
	err = col.FindOne(context.Background(), &Entity{
		Hi: "hello",
	}, &bbb)
	assert.NoError(t, err)

}
