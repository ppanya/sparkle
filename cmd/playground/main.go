package main

import (
	"context"
	"fmt"
	"github.com/octofoxio/foundation"
	"github.com/octofoxio/sparkle"
	"github.com/octofoxio/sparkle/external/mongodb"
	"github.com/octofoxio/sparkle/internal/migrate"
	commonv1 "github.com/octofoxio/sparkle/pkg/common/v1"
	sparklecrypto "github.com/octofoxio/sparkle/pkg/crypto"
	entitiesv1 "github.com/octofoxio/sparkle/pkg/entities/v1"
	"github.com/octofoxio/sparkle/pkg/svcs"
	svcsv1 "github.com/octofoxio/sparkle/pkg/svcs/v1"
	"github.com/octofoxio/sparkle/pkg/testutils"
	"math"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"
)

var config *sparkle.Config

func StartPlaygroundService() {
	wd, _ := os.Getwd()
	system := foundation.NewFileSystem(path.Join(wd, "./resources"), foundation.StaticMode_LOCAL)
	db := mongodb.NewLocal(testutils.DatabaseName)
	config = sparkle.NewConfig(system).
		SetDatabase(db).
		SetAddress(sparkle.LocalSparkleServiceURL).
		SetHost(sparkle.LocalHostURL).
		SetTokenSigner(sparklecrypto.NewBase64TokenSigner()).
		// Set default email template to blank for easy to testing
		SetDefaultEmailTemplate("{{.ConfirmUrl}}")

	var err error
	err = migrate.DropMongoCollection(db, config)
	if err != nil {
		panic(err)
	}
	err = migrate.MigrateMongoCollection(db, config)
	if err != nil {
		panic(err)
	}

	grpcServer, httpHandler := svcs.NewSparkleV1(config)
	c := context.Background()
	go sparkle.MustListenAndServeHTTP(httpHandler, config.Host.Port())
	go sparkle.MustListenAndServeTCP(grpcServer, config.Address.Port())
	<-c.Done()
}

func main() {

	go StartPlaygroundService()

	time.Sleep(time.Second * 2)
	conn := foundation.MakeDialOrPanic(config.Address.Host)
	client := svcsv1.NewSparkleClient(conn)

	output, err := client.Register(context.Background(),
		&svcsv1.RegisterInput{
			RegisterInputData: &svcsv1.RegisterInput_Email{
				Email: &svcsv1.RegisterWithEmailInput{
					Email:         commonv1.NotNullString("johnny@appleseed.com"),
					DisplayName:   commonv1.NotNullString("John"),
					CallbackURL:   commonv1.NotNullString("https://www.google.com"),
					PlainPassword: commonv1.NotNullString("something"),
					PhoneNumber:   commonv1.NotNullString("+66875969139"),
					FullName:      commonv1.NotNullString("Johnny Apple Seed"),
					Gender:        commonv1.Gender_Female,
				},
			},
		},
	)

	if err != nil {
		panic(err)
	}
	fmt.Println(output.Result.GetID().GetData())
	confirmationURL, err := url.Parse(config.EmailSender.(*sparkle.ConsoleEmailSender).Inbox[0][3])
	if err != nil {
		panic(err)
	}
	fmt.Printf("confirmation url: %s\n", confirmationURL)
	resp, err := http.Get(confirmationURL.String())
	if err != nil {
		panic(err)
	}

	fmt.Println(resp.Status)
	//now := time.Now()
	//guard := monkey.Patch(time.Now, func() time.Time {
	//	// because in setup_test.go
	//	// we use accessTokenLifetime = 1 minute
	//	return now.Add(time.Hour * 3600)
	//})
	//defer guard.Restore()
	time.Sleep(time.Second * 1)
	token := commonv1.NotNullString(confirmationURL.Query().Get("token"))
	validateOutput, err := client.ValidateAccessToken(context.Background(), &svcsv1.ValidateAccessTokenInput{
		AccessToken: token,
	})

	if err != nil {
		panic(err)
	}

	sessionContext := foundation.AppendAuthorizationToContext(context.Background(), token.GetData())
	userprofile, err := client.GetMyProfile(sessionContext, &svcsv1.GetMyProfileInput{
		SiteName: commonv1.NotNullString("default"),
	})

	if err != nil {
		panic(err)
	}

	fmt.Println(userprofile.String())
	identity, err := client.GetIdentity(sessionContext, &svcsv1.GetIdentityInput{
		SiteName: commonv1.NotNullString("reeeed"),
	})
	if err != nil {
		panic(err)
	}

	identityPutOutput, err := client.PutIdentity(sessionContext, &svcsv1.PutIdentityInput{
		SiteName: commonv1.NotNullString("reeeed"),
		Data: &entitiesv1.Identity{
			DisplayName: commonv1.NotNullString("James"),
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("===== token validate =====")
	fmt.Println(validateOutput.String())
	fmt.Println("===== profile =====")
	fmt.Println(userprofile.Result.String())

	fmt.Println("===== identity default =====")
	fmt.Println(identity.Result.String())
	fmt.Println("===== identity for reeeed =====")
	fmt.Println(identityPutOutput.Result.String())

	loginOutput, err := client.Login(context.Background(), &svcsv1.LoginInput{
		LoginInputData: &svcsv1.LoginInput_Email{
			Email: &svcsv1.LoginInputWithEmail{
				Email:         userprofile.Result.GetEmail(),
				PlainPassword: commonv1.NotNullString("something"),
			},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(loginOutput.Result.String())

	var (
		channelID     = foundation.EnvString("SPARKLE_LINE_CHANNEL_ID", "")
		channelSecret = foundation.EnvString("SPARKLE_LINE_SECRET", "")
	)
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
