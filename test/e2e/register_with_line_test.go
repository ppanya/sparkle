package e2e

import (
	"context"
	"testing"
	"time"

	"github.com/octofoxio/foundation"

	"bou.ke/monkey"
	"github.com/octofoxio/sparkle"
	commonv1 "github.com/octofoxio/sparkle/pkg/common/v1"
	entitiesv1 "github.com/octofoxio/sparkle/pkg/entities/v1"
	svcsv1 "github.com/octofoxio/sparkle/pkg/svcs/v1"
	"github.com/octofoxio/sparkle/pkg/testutils"
	"github.com/stretchr/testify/assert"
)

var (
	// TODO Remove SECRET
	channelID     = foundation.EnvString("SPARKLE_LINE_CHANNEL_ID", "")
	channelSecret = foundation.EnvString("SPARKLE_LINE_SECRET", "")
)

func Test_E2E_RegisterWithLine(t *testing.T) {

	if testing.Short() {
		t.Skip()
	}
	testutils.NewSuite(t, func(t *testing.T, database sparkle.Database, clients *testutils.SuiteClients) {
		ctx := context.Background()

		var tokenCh = make(chan testutils.LineSession)
		go testutils.LineLogin(channelID, channelSecret, tokenCh)
		lineSession := <-tokenCh

		registerWithLineOutput, err := clients.Sparkle.Register(context.Background(), &svcsv1.RegisterInput{
			RegisterInputData: &svcsv1.RegisterInput_Line{
				Line: &svcsv1.RegisterWithLineInput{
					AccessToken: commonv1.NotNullString(lineSession.AccessToken),
				},
			},
		})
		if err != nil {
			panic(err)
		}
		assert.NoError(t, err)
		assert.NotEmpty(t, registerWithLineOutput)

		loginWithLineOutput, err := clients.Sparkle.Login(context.Background(), &svcsv1.LoginInput{
			LoginInputData: &svcsv1.LoginInput_Line{
				Line: &svcsv1.LoginInputWithLine{
					LineAccessToken: commonv1.NotNullString(lineSession.AccessToken),
				},
			},
		})

		userID := loginWithLineOutput.GetResult().GetUserID()
		sessionCollections := database.Collection(sparkle.SessionCollectionName)
		var session entitiesv1.Session
		err = sessionCollections.FindOne(ctx, &entitiesv1.Session{UserID: userID}, &session)
		assert.NoError(t, err)
		assert.Equal(t, session.GetUserID(), userID)
		assert.NotNil(t, session.GetAccessToken())

		t.Run("identity is created", func(t *testing.T) {
			identityCollections := database.Collection(sparkle.IdentityCollectionName)
			var identity entitiesv1.IdentityRecord
			err := identityCollections.FindOne(ctx, &entitiesv1.IdentityRecord{UserID: userID.GetData()}, &identity)
			assert.NoError(t, err)
			assert.NotNil(t, identity)
		})

		t.Run("access token is valid after register and confirm email", func(t *testing.T) {
			validateOutput, err := clients.Sparkle.ValidateAccessToken(context.Background(), &svcsv1.ValidateAccessTokenInput{
				AccessToken: session.GetAccessToken(),
			})
			assert.NoError(t, err)
			assert.EqualValues(t, validateOutput.GetResult().GetIsValid(), true)
		})

		now := time.Now()
		guard := monkey.Patch(time.Now, func() time.Time {
			// because in setup_test.go
			// we use accessTokenLifetime = 1 minute
			return now.Add(time.Hour * 3600)
		})
		t.Run("access token should invalid after not use for a period of time", func(t *testing.T) {
			validateOutput, err := clients.Sparkle.ValidateAccessToken(context.Background(), &svcsv1.ValidateAccessTokenInput{
				AccessToken: session.GetAccessToken(),
			})
			assert.NoError(t, err)
			assert.EqualValues(t, validateOutput.GetResult().GetIsValid(), false)
		})
		defer guard.Restore()
	})

}
