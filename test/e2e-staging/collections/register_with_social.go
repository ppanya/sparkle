package collections

import (
	"fmt"
	"github.com/octofoxio/foundation"
	"github.com/octofoxio/sparkle/pkg/common/v1"
	"github.com/octofoxio/sparkle/pkg/svcs/v1"
	"github.com/octofoxio/sparkle/pkg/testutils"
	"golang.org/x/net/context"
	"math"
	"time"
)

var (
	channelID     = foundation.EnvString("SPARKLE_LINE_CHANNEL_ID", "1597323211")
	channelSecret = foundation.EnvString("SPARKLE_LINE_SECRET", "4c8cbc5d3d82b36a754a3b71017abaa8")
)

func RegisterWithLine(client svcsv1.SparkleClient) {
	fmt.Println("====== Begin LINE register =====")
	var tokenCh = make(chan testutils.LineSession)
	go testutils.LineLogin(channelID, channelSecret, tokenCh)
	var lineSession = <-tokenCh

	registerWithLineOutput, err := client.Register(context.Background(), &svcsv1.RegisterInput{
		RegisterInputData: &svcsv1.RegisterInput_Line{
			Line: &svcsv1.RegisterWithLineInput{
				AccessToken: commonv1.NotNullString(lineSession.AccessToken),
			},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("======= line register output =======")
	fmt.Println(registerWithLineOutput.String())

	loginWithLineOutput, err := client.Login(context.Background(), &svcsv1.LoginInput{
		LoginInputData: &svcsv1.LoginInput_Line{
			Line: &svcsv1.LoginInputWithLine{
				LineAccessToken: commonv1.NotNullString(lineSession.AccessToken),
			},
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("===== line login output =====")
	fmt.Println(loginWithLineOutput.String())

	// call confirm email
	time.Sleep(math.MaxInt64)
}
