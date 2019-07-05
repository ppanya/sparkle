package rand

import (
	"github.com/brianvoe/gofakeit"
	commonv1 "github.com/octofoxio/sparkle/pkg/common/v1"
	entitiesv1 "github.com/octofoxio/sparkle/pkg/entities/v1"
)

func User() *entitiesv1.User {
	return &entitiesv1.User{
		Email:  commonv1.NotNullString(gofakeit.Email()),
		Status: entitiesv1.UserStatus_WaitingForEmailVerification,
	}
}
