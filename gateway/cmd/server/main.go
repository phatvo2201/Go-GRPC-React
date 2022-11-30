package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	userinfo "github.com/phatbb/userinfo/pb"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
)

var (
	grpcUserServerEndpoint = flag.String("grpc-user-server-endpoint", "host.docker.internal:9091", "gRPC user server endpoint")
	grpcAuthServerEndpoint = flag.String("grpc-auth-server-endpoint", "host.docker.internal:9090", "gRPC auth server endpoint")
)

//custom to disable grpc-metadata prefix

func CustomMatcher(key string) (string, bool) {
	return fmt.Sprintf("%s", key), true

}

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Register gRPC server endpoint
	// new server mux with options WithOutgoingHeaderMatcher
	mux := runtime.NewServeMux(runtime.WithOutgoingHeaderMatcher(CustomMatcher))

	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{
			http.MethodPost,
			http.MethodGet,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	handler := cors.Handler(mux)

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := userinfo.RegisterUserServiceHandlerFromEndpoint(ctx, mux, *grpcUserServerEndpoint, opts)
	log.Println("regis gate way with user handler rpc")
	if err != nil {
		log.Fatal("can not regis user gw")
	}

	err = userinfo.RegisterAuthenServiceHandlerFromEndpoint(ctx, mux, *grpcAuthServerEndpoint, opts)
	log.Println("regis gate way with auth handler rpc")
	if err != nil {
		log.Fatal("can not regis gw")
	}

	return http.ListenAndServe(":8080", handler)
}
func main() {
	if err := run(); err != nil {
		glog.Fatal(err)
	}

}
