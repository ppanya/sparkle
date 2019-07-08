package sparkleuc

import (
	"context"
	"errors"
	entitiesv1 "github.com/octofoxio/sparkle/pkg/entities/v1"
	svcsv1 "github.com/octofoxio/sparkle/pkg/svcs/v1"
	"go.mongodb.org/mongo-driver/mongo"
)

func (l *LoginUseCase) LoginWithEmail(ctx context.Context, in *svcsv1.LoginInput_LoginInputWithEmail) (*entitiesv1.SessionRecord, error) {
	user, err := l.user.FindOne(ctx, &entitiesv1.User{
		Email: in.Email,
	})

	if err != nil && err == mongo.ErrNoDocuments {
		return nil, errors.New("invalid user credential")
	} else if err != nil {
		return nil, err
	}

	if !user.ValidatePassword(in.PlainPassword.GetData()) {
		return nil, errors.New("invalid user credential")
	}

	s, err := l.CreateSession(ctx, user.ID.GetData())
	if err != nil {
		return nil, err
	}

	return s, nil

}
