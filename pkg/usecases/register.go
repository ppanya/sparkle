package sparkleuc

import (
	"context"
	entitiesv1 "github.com/octofoxio/sparkle/pkg/entities/v1"
	svcsv1 "github.com/octofoxio/sparkle/pkg/svcs/v1"
)

func (s *SparkleUseCase) RegisterWithEmail(c context.Context, in *svcsv1.RegisterWithEmailInput) (*svcsv1.RegisterWithEmailOutput, error) {

	ID, err := s.user.Create(c, &entitiesv1.UserRecord{
		User: entitiesv1.User{
			Status: entitiesv1.UserStatus_WaitingForEmailVerification,
			Email:  in.Email,
		},
	})

}
