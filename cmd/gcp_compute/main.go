package main

import (
	"context"
	"fmt"
	"github.com/octofoxio/foundation"
	"github.com/octofoxio/foundation/logger"
	"github.com/octofoxio/sparkle"
	_ "github.com/octofoxio/sparkle/cmd/statik/statik"
	"github.com/octofoxio/sparkle/external/mailgun"
	"github.com/octofoxio/sparkle/external/mongodb"
	"github.com/octofoxio/sparkle/internal/migrate"
	"github.com/octofoxio/sparkle/pkg/crypto"
	"github.com/octofoxio/sparkle/pkg/svcs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net"
	"net/http"
)

func main() {

	log := logger.New("sparkle-gcp-compute")

	var (
		MongoDBURL          = foundation.EnvStringOrPanic("SPARKLE_MONGODB_URL")
		Host                = foundation.EnvStringOrPanic("SPARKLE_HOST")
		ServiceAddress      = foundation.EnvStringOrPanic("SPARKLE_SERVICE_ADDRESS")
		mailGunDomain       = foundation.EnvStringOrPanic("SPARKLE_MAILGUN_DOMAIN")
		mailGunApiKey       = foundation.EnvStringOrPanic("SPARKLE_MAILGUN_APIKEY")
		mailGunEmailAddress = foundation.EnvStringOrPanic("SPARKLE_MAILGUN_EMAIL")
	)
	system := foundation.NewFileSystem("", foundation.StaticMode_Statik)
	client, err := mongo.NewClient(
		options.Client().ApplyURI(MongoDBURL))
	if err != nil {
		panic(err)
	}
	err = client.Connect(context.Background())
	if err != nil {
		panic(err)
	}
	db := mongodb.New(client.Database("sparkle"))

	var config = sparkle.NewConfig(system).
		SetEmailSender(mailgun.NewMailGunEmailSender(mailGunDomain, mailGunApiKey)).
		SetDefaultEmailSenderAddress(mailGunEmailAddress).
		SetDatabase(db).
		SetHost(Host).
		SetAddress(ServiceAddress).
		SetTokenSigner(sparklecrypto.NewBase64TokenSigner())
	migrate.MustMigrateMongoCollection(db, config)

	GRPCServer, HTTPServer := svcs.NewSparkleV1(config)
	go func() {
		log.Printf("grpc start! (%s)", config.Address.String())
		log.Printf("config.Address.Port(): %s", config.Address.Port())
		lis, _ := net.Listen("tcp", fmt.Sprintf(":%s", config.Address.Port()))
		err := GRPCServer.Serve(lis)
		if err != nil {
			panic(err)
		}
	}()

	log.Printf("http start! (%s)", config.Host.String())
	log.Printf("config.Host.Port(): %s", config.Host.Port())
	err = http.ListenAndServe(fmt.Sprintf(":%s", config.Host.Port()), HTTPServer)
	if err != nil {
		panic(err)
	}

}
