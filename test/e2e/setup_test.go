package e2e

import (
	"github.com/octofoxio/foundation"
	"github.com/octofoxio/sparkle"
	"github.com/octofoxio/sparkle/external/mongodb"
	"github.com/octofoxio/sparkle/internal/migrate"
	"github.com/octofoxio/sparkle/pkg/svcs"
	"github.com/octofoxio/sparkle/pkg/testutils"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"path"
	"testing"
)

var (
	config *sparkle.Config
)

func TestMain(m *testing.M) {

	wd, _ := os.Getwd()
	system := foundation.NewFileSystem(path.Join(wd, "../../resources"), foundation.StaticMode_LOCAL)
	db := mongodb.NewLocal(testutils.DatabaseName)

	config = sparkle.NewConfig(system).
		SetDatabase(db).
		SetHost(sparkle.LocalHostURL).
		SetAddress(sparkle.LocalSparkleServiceURL).
		UseJWTSignerWithHMAC256("integration-test")

	serv, httpserv := svcs.NewSparkleV1(config)
	go func(s *grpc.Server) {
		lis, _ := net.Listen("tcp", config.Address.String())
		_ = s.Serve(lis)
	}(serv)

	go func(h http.Handler) {
		err := http.ListenAndServe(":"+config.Host.Opaque, h)
		if err != nil {
			panic(err)
		}
	}(httpserv)

	err := migrate.MigrateMongoCollection(db, config)
	if err != nil {
		panic(err)
	}

	c := m.Run()
	os.Exit(c)

}
