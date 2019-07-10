package line

import (
	"context"
	"fmt"
	"github.com/octofoxio/foundation"
	"github.com/octofoxio/sparkle/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	channelID     = foundation.EnvString("SPARKLE_LINE_CHANNEL_ID", "")
	channelSecret = foundation.EnvString("SPARKLE_LINE_SECRET", "")
)

func TestGetProfile(t *testing.T) {

	if testing.Short() || len(channelID) == 0 || len(channelSecret) == 0 {
		fmt.Println("skipping LINE integration test, please provide SPARKLE_LINE_CHANNEL_ID and SPARKLE_LINE_SECRET")
		t.Skip()
	}

	var resultChannel = make(chan testutils.LineSession)
	go testutils.LineLogin(channelID, channelSecret, resultChannel)
	session := <-resultChannel

	t.Logf("access token: %s\n", session.AccessToken)

	cc := DefaultLineClient{}
	result, err := cc.GetProfile(context.Background(), session.AccessToken)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.DisplayName)
	assert.NotEmpty(t, result.PictureURL)
	t.Logf("display name: %s", result.DisplayName)
	t.Logf("picture URL: %s", result.PictureURL)

}
