package main

import (
	"github.com/octofoxio/foundation"
	"github.com/octofoxio/sparkle/pkg/svcs/v1"
	"github.com/octofoxio/sparkle/test/e2e-staging/collections"
	"net/url"
)

var (
	sparkleHTTP, _ = url.Parse("http://35.198.223.236:3009")
	sparkleGRPC, _ = url.Parse("//35.198.223.236:3019")
)

func main() {
	conn := foundation.MakeDialOrPanic(sparkleGRPC.Host)
	client := svcsv1.NewSparkleClient(conn)

	collections.RegisterWithEmail(client)
	//collections.RegisterWithLine(client)
}
