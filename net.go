package sparkle

import (
	"context"
	"fmt"
	entitiesv1 "github.com/octofoxio/sparkle/pkg/entities/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"net"
	"net/http"
)

var (
	AuthorizationContextKey = "authorization"
	UserRecordContextKey    = "user-record"
)

func AppendUserProfileToContext(ctx context.Context, record *entitiesv1.UserRecord) context.Context {
	return context.WithValue(ctx, UserRecordContextKey, record)
}
func GetUserProfileFromContext(ctx context.Context) (out *entitiesv1.UserRecord, ok bool) {
	if ctx.Value(UserRecordContextKey) == nil {
		return nil, false
	}
	if u, ok := ctx.Value(UserRecordContextKey).(*entitiesv1.UserRecord); ok {
		return u, true
	} else {
		return nil, false
	}
}

func AppendAccessTokenToOutgoingContext(ctx context.Context, accessToken string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, AuthorizationContextKey, accessToken)
}

func GetAccessTokenFromIncomingContext(ctx context.Context) (string, bool) {
	md, isOk := metadata.FromIncomingContext(ctx)
	if !isOk {
		return "", false
	}

	data := md.Get(AuthorizationContextKey)
	if len(data) == 0 {
		return "", false
	}
	return data[0], true
}

func MustListenAndServeTCP(server *grpc.Server, port string) {
	fmt.Printf("====== TCP service start at port %s ======\n", port)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		panic(err)
	}
	_ = server.Serve(lis)
}

func MustListenAndServeHTTP(handler http.Handler, port string) {
	fmt.Printf("====== HTTP service start at port %s ======\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), handler)
	if err != nil {
		panic(err)
	}
}
