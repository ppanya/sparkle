package testutils

import (
	"github.com/octofoxio/foundation"
	"github.com/octofoxio/sparkle"
	"github.com/octofoxio/sparkle/external/mongodb"
	svcsv1 "github.com/octofoxio/sparkle/pkg/svcs/v1"
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
	clients := &SuiteClients{
		Sparkle: svcsv1.NewSparkleClient(foundation.MakeDialOrPanic(sparkle.LocalSparkleServiceURL)),
		Spike:   svcsv1.NewSpikeClient(foundation.MakeDialOrPanic(sparkle.LocalSpikeServiceURL)),
	}
	db := mongodb.NewLocal(DatabaseName)
	fn(t, db, clients)
}
