package sparkleuc

import (
	"context"
	"errors"
	commonv1 "github.com/octofoxio/sparkle/pkg/common/v1"
	sparklecrypto "github.com/octofoxio/sparkle/pkg/crypto"
	entitiesv1 "github.com/octofoxio/sparkle/pkg/entities/v1"
	sparklerepo "github.com/octofoxio/sparkle/pkg/repositories"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type LoginUseCase struct {
	session  sparklerepo.SessionRepository
	identity sparklerepo.IdentityRepository
	user     sparklerepo.UserRepository
	signer   sparklecrypto.TokenSigner
}

func NewLoginUseCase(session sparklerepo.SessionRepository, identity sparklerepo.IdentityRepository, user sparklerepo.UserRepository, signer sparklecrypto.TokenSigner) *LoginUseCase {
	return &LoginUseCase{session: session, identity: identity, user: user, signer: signer}
}

func (l *LoginUseCase) ValidateSession(ctx context.Context, accessToken string) (*entitiesv1.UserRecord, error) {

	session, err := l.session.FindOne(ctx, &entitiesv1.Session{
		AccessToken: commonv1.NotNullString(accessToken),
	})

	if err == mongo.ErrNoDocuments {
		return nil, errors.New("session not found")
	}

	if err != nil {
		return nil, err
	}

	if err := session.IsValid(); err != nil {
		return nil, err
	}

	user, err := l.user.FindByID(ctx, session.UserID.GetData())
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("user not found, maybe invalid session")
	}
	return user, err

}

func (l *LoginUseCase) CreateSession(ctx context.Context, userID string) (*entitiesv1.SessionRecord, error) {
	accessToken, err := NewSession(l.signer, userID)
	session := &entitiesv1.SessionRecord{
		Session: entitiesv1.Session{
			UserID:          commonv1.NotNullString(userID),
			CreatedAt:       commonv1.NewTimestamp(time.Now()),
			AccessToken:     commonv1.NotNullString(accessToken),
			LatestVisitedAt: commonv1.NewTimestamp(time.Now()),
		},
	}
	_, err = l.session.Create(ctx, session)
	if err != nil {
		return nil, err
	}
	return session, err
}
