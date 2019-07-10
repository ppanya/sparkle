package e2e

import (
	"bou.ke/monkey"
	"context"
	"github.com/octofoxio/sparkle"
	commonv1 "github.com/octofoxio/sparkle/pkg/common/v1"
	"github.com/octofoxio/sparkle/pkg/rand"
	svcsv1 "github.com/octofoxio/sparkle/pkg/svcs/v1"
	"github.com/octofoxio/sparkle/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestRegisterWithEmail(t *testing.T) {

	if testing.Short() {
		t.Skip()
	}
	testutils.NewSuite(t, func(t *testing.T, database sparkle.Database, clients *testutils.SuiteClients) {

		registerInputData := rand.RegisterWithEmailInput()
		registerInputData.CallbackURL = commonv1.NotNullString("https://www.google.com")
		output, err := clients.Sparkle.Register(context.Background(), &svcsv1.RegisterInput{
			RegisterInputData: &svcsv1.RegisterInput_Email{
				Email: registerInputData,
			},
		})
		assert.NoError(t, err)
		assert.NotEmpty(t, output)

		mail := config.EmailSender.(*sparkle.ConsoleEmailSender).Inbox[0][2]
		assert.NotEmpty(t, mail)
		resp, err := http.Get(mail)
		assert.NoError(t, err)
		assert.Equal(t, resp.StatusCode, http.StatusOK)

		confirmationURL, err := url.Parse(mail)
		assert.NoError(t, err)
		token := commonv1.NotNullString(confirmationURL.Query().Get("token"))

		t.Run("access token is valid after register and confirm email", func(t *testing.T) {
			validateOutput, err := clients.Sparkle.ValidateAccessToken(context.Background(), &svcsv1.ValidateAccessTokenInput{
				AccessToken: token,
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
				AccessToken: token,
			})
			assert.NoError(t, err)
			assert.EqualValues(t, validateOutput.GetResult().GetIsValid(), false)
		})
		defer guard.Restore()

	})

}
