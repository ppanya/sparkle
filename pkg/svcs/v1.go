package svcs

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/octofoxio/foundation"
	"github.com/octofoxio/sparkle"
	"github.com/octofoxio/sparkle/pkg/common"
	"github.com/octofoxio/sparkle/pkg/endpoints"
	"github.com/octofoxio/sparkle/pkg/repositories"
	"github.com/octofoxio/sparkle/pkg/svcs/v1"
	"github.com/octofoxio/sparkle/pkg/usecases"
	"google.golang.org/grpc"
	"net/http"
	"runtime/debug"
)

func NewSpikeV1(config *sparkle.Config) *grpc.Server {
	var (
		spikeEndpoint = endpoints.NewSpikeEndpoints()
		serv          = foundation.NewGRPCServer(
			common.NewConfigLoaderInterceptor(
				config,
			),
		)
	)
	svcsv1.RegisterSpikeServer(serv, spikeEndpoint)
	return serv
}

func NewSparkleV1(config *sparkle.Config) (*grpc.Server, http.Handler) {
	var (
		identityRepository = sparklerepo.NewDefaultIdentityRepository(config.Database.Collection(config.IdentityCollectionName))
		userRepository     = sparklerepo.NewDefaultUserRepository(config.Database.Collection(config.UserCollectionName))
		sessionRepository  = sparklerepo.NewDefaultSessionRepository(config.Database.Collection(config.SessionCollectionName))

		registerUseCase = sparkleuc.NewRegisterUseCase(config.TokenSigner, sessionRepository, identityRepository, userRepository, config.EmailSender, config.Fs)
		loginUseCase    = sparkleuc.NewLoginUseCase(sessionRepository, identityRepository, userRepository, config.TokenSigner)

		sparkleEndpoint = endpoints.NewSparkleEndpoints(registerUseCase, loginUseCase)

		serv = foundation.NewGRPCServer(
			common.NewConfigLoaderInterceptor(
				config,
			),
			grpc_recovery.UnaryServerInterceptor(
				grpc_recovery.WithRecoveryHandler(func(p interface{}) (err error) {

					// recovery panic from handler
					fmt.Println(string(debug.Stack()))
					switch b := p.(type) {
					case error:
						return b
					case string:
						return errors.New(b)
					default:
						fmt.Println("==== UNKNOWN ERROR OCCUR BEGIN ====")
						fmt.Println(p)
						fmt.Println("==== UNKNOWN ERROR OCCUR END ====")
						return errors.New("unknown error")
					}

				}),
			),
			func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

				// extract user session from context
				var accessToken, ok = sparkle.GetAccessTokenFromIncomingContext(ctx)
				if ok {
					o, err := loginUseCase.ValidateSession(ctx, accessToken)
					if err != nil {
						return nil, err
					}
					ctx = sparkle.AppendUserProfileToContext(ctx, o)
					return handler(ctx, req)
				} else {
					return handler(ctx, req)
				}
			},
		)
	)
	r := NewSparkleViews(config, sparkleEndpoint)
	svcsv1.RegisterSparkleServer(serv, sparkleEndpoint)

	if config.TokenSigner == nil {
		panic("token signer is nil")
	}

	return serv, r
}

func ConfirmEmailHTTPHandler(config *sparkle.Config, e *endpoints.SparkleEndpoints) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			c           = common.AppendConfigToContext(context.Background(), config)
			accessToken = r.FormValue("token")
			callbackURL = r.FormValue("callbackURL")
		)

		err := e.ConfirmEmailByAccessTokenHandler(c, accessToken, callbackURL)
		if err != nil {
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusBadRequest)
			_, err = w.Write([]byte("something wrong: " + err.Error()))
			return
		} else {
			http.Redirect(w, r, callbackURL, http.StatusTemporaryRedirect)
			return
		}
	}
}

func NewSparkleViews(config *sparkle.Config, endpoint *endpoints.SparkleEndpoints) http.Handler {
	r := mux.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	})
	r.Path("/c").Methods("get").HandlerFunc(ConfirmEmailHTTPHandler(config, endpoint))
	return r
}
