package testutils

import (
	"context"
	"net/url"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/octofoxio/foundation"
	"github.com/octofoxio/sparkle"
	"github.com/octofoxio/sparkle/external/mongodb"
	svcsv1 "github.com/octofoxio/sparkle/pkg/svcs/v1"
)

type SuiteClients struct {
	Sparkle svcsv1.SparkleClient
	Spike   svcsv1.SpikeClient
}

var (
	DatabaseName = foundation.EnvString("SPARKLE_MONGODB_NAME", "sparkle-test")
	HTTPEndpoint = foundation.EnvString("SPARKLE_HTTP_ENDPOINT", sparkle.LocalHostURL)
	GRPCEndpoint = foundation.EnvString("SPARKLE_GRPC_ENDPOINT", sparkle.LocalSparkleServiceURL)
	e2e          = foundation.EnvString("SPARKLE_E2E_TEST", "")
	DatabaseURL  = foundation.EnvString("SPARKLE_MONGODB_URL", sparkle.LocalMongoDBURL)
)

// NewSuite will create anything that
// need to use at test suite
// this suite require mongodb run in replica set
func NewSuite(t *testing.T, fn func(t *testing.T, database sparkle.Database, clients *SuiteClients)) {
	sparkURL, _ := url.Parse(sparkle.LocalSparkleServiceURL)
	spikeURL, _ := url.Parse(sparkle.LocalSpikeServiceURL)

	clients := &SuiteClients{
		Sparkle: svcsv1.NewSparkleClient(foundation.MakeDialOrPanic(sparkURL.Host)),
		Spike:   svcsv1.NewSpikeClient(foundation.MakeDialOrPanic(spikeURL.Host)),
	}
	var db *mongodb.MongoDatabase
	if e2e == "1" {
		client, err := mongo.NewClient(
			options.Client().ApplyURI(DatabaseURL))
		if err != nil {
			panic(err)
		}
		err = client.Connect(context.Background())
		if err != nil {
			panic(err)
		}
		db = mongodb.New(client.Database(DatabaseName))
	} else {
		db = mongodb.NewLocal(DatabaseName)
	}
	fn(t, db, clients)
}
