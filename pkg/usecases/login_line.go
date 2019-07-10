package sparkleuc

import (
	"context"
	"errors"
	"github.com/octofoxio/sparkle"
	commonv1 "github.com/octofoxio/sparkle/pkg/common/v1"
	entitiesv1 "github.com/octofoxio/sparkle/pkg/entities/v1"
	svcsv1 "github.com/octofoxio/sparkle/pkg/svcs/v1"
)

func (l *LoginUseCase) LoginWithLine(ctx context.Context, in *svcsv1.LoginInputWithLine) (*entitiesv1.SessionRecord, error) {

	lineProfile, err := l.lineClient.GetProfile(ctx, in.LineAccessToken.GetData())
	if err != nil {
		return nil, err
	}

	user, err := l.user.FindOne(ctx, &entitiesv1.UserRecord{
		LineID: commonv1.NotNullString(lineProfile.UserID),
	})

	if err != nil && err == sparkle.ErrNotFound {
		return nil, errors.New("invalid user credential")
	} else if err != nil {
		return nil, err
	}

	s, err := l.CreateSession(ctx, user.ID.GetData())
	if err != nil {
		return nil, err
	}

	return s, nil

}
