package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	userinfo "github.com/phatbb/userinfo/pb"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	grpcUserServerEndpoint = flag.String("grpc-user-server-endpoint", "192.168.5.33:9091", "gRPC user server endpoint")
	grpcAuthServerEndpoint = flag.String("grpc-auth-server-endpoint", "192.168.5.33:9090", "gRPC auth server endpoint")
)

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
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
		log.Fatal("can not regis gw")
	}

	err = userinfo.RegisterSimpleBankHandlerFromEndpoint(ctx, mux, *grpcAuthServerEndpoint, opts)
	log.Println("regis gate way with auth handler rpc")
	if err != nil {
		log.Fatal("can not regis gw")
	}
	// err = userinfo.(ctx, mux, *grpcServerEndpoint, opts)
	// log.Println("regis gate way with user handler rpc")
	// if err != nil {
	// 	log.Fatal("can not regis gw")
	// }

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	return http.ListenAndServe(":8080", handler)
}
func main() {
	if err := run(); err != nil {
		glog.Fatal(err)
	}

}
