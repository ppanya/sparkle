package main

import (
	"context"
	"fmt"
	"github.com/octofoxio/foundation"
	"github.com/octofoxio/sparkle"
	"github.com/octofoxio/sparkle/external/mongodb"
	commonv1 "github.com/octofoxio/sparkle/pkg/common/v1"
	sparklecrypto "github.com/octofoxio/sparkle/pkg/crypto"
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

	grpcServer, httpHandler := svcs.NewSparkleV1(config)
	c := context.Background()
	go sparkle.MustListenAndServeHTTP(httpHandler, config.Host.Host)
	go sparkle.MustListenAndServeTCP(grpcServer, config.Address.String())
	<-c.Done()
}

func main() {

	go StartPlaygroundService()

	time.Sleep(time.Second * 2)
	conn := foundation.MakeDialOrPanic(sparkle.LocalSparkleServiceURL)
	client := svcsv1.NewSparkleClient(conn)

	output, err := client.RegisterWithEmail(context.Background(), &svcsv1.RegisterWithEmailInput{
		Email:         commonv1.NotNullString("johnny@appleseed.com"),
		DisplayName:   commonv1.NotNullString(""),
		CallbackURL:   commonv1.NotNullString("https://www.google.com"),
		PlainPassword: commonv1.NotNullString("something"),
		PhoneNumber:   commonv1.NotNullString("+66875969139"),
		FullName:      commonv1.NotNullString("Johnny Apple Seed"),
		Gender:        commonv1.Gender_Female,
	})

	if err != nil {
		panic(err)
	}
	fmt.Println(output.Result.GetID().GetData())
	confirmationURL, err := url.Parse(config.EmailSender.(*sparkle.ConsoleEmailSender).Inbox[0][2])
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
	fmt.Println(validateOutput.String())

	sessionContext := foundation.AppendAuthorizationToContext(context.Background(), token.GetData())
	userprofile, err := client.GetProfile(sessionContext, &svcsv1.GetProfileInput{
		SiteName: commonv1.NotNullString("default"),
	})

	// call confirm email
	time.Sleep(math.MaxInt64)

}
