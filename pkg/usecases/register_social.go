package sparkleuc

import (
	"context"
	"errors"
	commonv1 "github.com/octofoxio/sparkle/pkg/common/v1"
	entitiesv1 "github.com/octofoxio/sparkle/pkg/entities/v1"
	svcsv1 "github.com/octofoxio/sparkle/pkg/svcs/v1"
	"time"
)

func (s *RegisterUseCase) RegisterWithLine(c context.Context, in *svcsv1.RegisterWithLineInput) (out *entitiesv1.User, err error) {

	lineProfile, err := s.lineClient.GetProfile(c, in.AccessToken.GetData())
	if err != nil {
		return nil, err
	}

	existsUser, err := s.user.FindOne(c, &entitiesv1.UserRecord{
		LineID: commonv1.NotNullString(lineProfile.UserID),
	})
	if existsUser != nil {
		return nil, errors.New(commonv1.T_SocialProfileHasAlreadyBeenUsed.String())
	}

	user := &entitiesv1.UserRecord{
		User: entitiesv1.User{
			Status:    entitiesv1.UserStatus_Active,
			CreatedAt: commonv1.NewTimestamp(time.Now()),
		},
		LineID:           commonv1.NotNullString(lineProfile.UserID),
		RegisterProvider: commonv1.RegisterProvider_LineProvider,
	}
	identity := &entitiesv1.IdentityRecord{
		SiteName: "default",
		UserID:   user.GetID().GetData(),
		Identity: entitiesv1.Identity{
			DisplayName:    commonv1.NotNullString(lineProfile.DisplayName),
			ProfilePicture: commonv1.NotNullString(lineProfile.PictureURL),
		},
	}

	ID, err := s.user.Create(c, user)
	_, err = s.identity.Create(c, identity)

	user.ID = commonv1.NotNullString(ID)
	return &user.User, nil
}
