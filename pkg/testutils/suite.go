package testutils

import (
	"github.com/octofoxio/foundation"
	"github.com/octofoxio/sparkle"
	"github.com/octofoxio/sparkle/external/mongodb"
	svcsv1 "github.com/octofoxio/sparkle/pkg/svcs/v1"
	"net/url"
	"testing"
)

type SuiteClients struct {
	Sparkle svcsv1.SparkleClient
	Spike   svcsv1.SpikeClient
}

var (
	DatabaseName = "sparkle-test"
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
	db := mongodb.NewLocal(DatabaseName)
	fn(t, db, clients)
}
