package sparklerepo

import (
	"context"
	"github.com/octofoxio/sparkle/external/mongodb"
	entitiesv1 "github.com/octofoxio/sparkle/pkg/entities/v1"
	"github.com/octofoxio/sparkle/pkg/rand"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultUserCreate(t *testing.T) {

	if testing.Short() {
		t.Skip()
	}

	mongoClient := mongodb.NewLocal("user-test")
	userRepository := NewDefaultUserRepository(mongoClient)

	ID, err := userRepository.Create(context.Background(), &entitiesv1.UserRecord{
		User: *rand.User(),
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, ID)

	createdUser, err := userRepository.FindByID(context.Background(), ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, createdUser)

	err = userRepository.UpdateByID(context.Background(), ID, &entitiesv1.UserRecord{
		User: entitiesv1.User{
			Status: entitiesv1.UserStatus_Active,
		},
	})
	assert.NoError(t, err)

	updatedUser, err := userRepository.FindByID(context.Background(), ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, updatedUser)
	assert.EqualValues(t, entitiesv1.UserStatus_Active, updatedUser.Status)

	_ = mongoClient.DB.Client().Disconnect(context.Background())

}
