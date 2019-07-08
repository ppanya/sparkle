package e2e

import (
	"context"
	"github.com/octofoxio/sparkle"
	"github.com/octofoxio/sparkle/pkg/rand"
	"github.com/octofoxio/sparkle/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegisterWithEmail(t *testing.T) {

	if testing.Short() {
		t.Skip()
	}
	testutils.NewSuite(t, func(t *testing.T, database sparkle.Database, clients *testutils.SuiteClients) {
		output, err := clients.Sparkle.RegisterWithEmail(context.Background(), rand.RegisterWithEmailInput())
		assert.NoError(t, err)
		assert.NotEmpty(t, output)

		mail := config.EmailSender.(*sparkle.ConsoleEmailSender).Inbox[0]
		assert.NotEmpty(t, mail)

	})

}
