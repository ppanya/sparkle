package sparkleuc

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	commonv1 "github.com/octofoxio/sparkle/pkg/common/v1"
	entitiesv1 "github.com/octofoxio/sparkle/pkg/entities/v1"
	"go.mongodb.org/mongo-driver/mongo"
)

func (d *RegisterUseCase) ConfirmEmailHandler(ctx aws.Context, accessToken string) error {

	session, err := d.session.FindOne(ctx, &entitiesv1.Session{
		AccessToken: commonv1.NotNullString(accessToken),
	})
	if err != nil {
		return err
	}
	fmt.Println(session.UserID.GetData())

	user, err := d.user.FindByID(ctx, session.UserID.GetData())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("user %s not found", session.UserID.GetData())
		}
		return err
	}

	if user.Status != entitiesv1.UserStatus_WaitingForEmailVerification {
		return errors.New("invalid user status, no confirmation need")
	}

	err = d.user.UpdateByID(ctx, session.UserID.GetData(), &entitiesv1.UserRecord{
		User: entitiesv1.User{
			Status: entitiesv1.UserStatus_Active,
		},
	})
	return err
}
