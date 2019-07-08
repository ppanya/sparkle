package common

import (
	"context"
	"github.com/octofoxio/sparkle"
	"google.golang.org/grpc"
)

var (
	ConfigContextKey = "config"
)

func NewConfigLoaderInterceptor(config *sparkle.Config) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		ctx = AppendConfigToContext(ctx, config)
		return handler(ctx, req)
	}
}

func AppendConfigToContext(ctx context.Context, config *sparkle.Config) context.Context {
	return context.WithValue(ctx, ConfigContextKey, config)
}

func GetConfigFromContext(ctx context.Context) *sparkle.Config {
	return ctx.Value(ConfigContextKey).(*sparkle.Config)
}
