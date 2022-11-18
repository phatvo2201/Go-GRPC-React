package main

import (
	"context"
	"flag"
	"log"
	"net"

	"github.com/phatbb/userinfo/config"
	"github.com/phatbb/userinfo/implgrpc"
	userinfo "github.com/phatbb/userinfo/pb"
	"github.com/phatbb/userinfo/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address = "192.168.5.33:9090"
)

func accessibleRole() map[string][]string {
	const infoservice = "/pb.UserService/"

	return map[string][]string{
		infoservice + "GetUserWalletInfo": {"user"},
		infoservice + "FindWalletByOwner": {"user"},
	}
}

func UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	log.Println("-------->???????????? unary intereptor:", info.FullMethod)
	return handler(ctx, req)

}

func startGrpcServer() {
	log.Println("start auth grpc server")
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Println("this is error from client side")
	}
	log.Println("connect  auth grpc server success")

	verifyClient := implgrpc.NewVerifyUserClient(conn)
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalln("cant not load config from env")
	}
	ctx := context.TODO()
	//create jwt manager
	mongoconn := options.Client().ApplyURI(config.DBUri)
	mongoClient, err := mongo.Connect(ctx, mongoconn)
	if err != nil {
		log.Fatal("cant not connect to the mongodb database")
	}
	usercollection := mongoClient.Database(config.DBName).Collection("users")

	walletCollection := mongoClient.Database(config.DBName).Collection("wallet")
	userService := service.NewUserService(usercollection, walletCollection, ctx)
	userServerHandle := implgrpc.NewUserServerImpl(config, userService, verifyClient)
	// jwtProvider := utils.NewJwtManager(config.AccessTokenPublicKey, config.AccessTokenPrivateKey, config.AccessTokenExpiresIn)

	grpcServer := grpc.NewServer()
	userinfo.RegisterUserServiceServer(grpcServer, userServerHandle)

	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatal("cannot create grpc server")
	}
	log.Printf("start grpc server on port %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot create grpc server")
	}

}

// func startGrpcGateway(listener net.Listener) error {
// 	mux := runtime.NewServeMux()
// 	ctx, cancel := context.WithCancel(context.Background())
// 	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

// 	grpcServerEndpoint := flag.String("grpc-server-endpoint", "localhost:8080", "gRPC server endpoint")
// 	err := userinfo.RegisterSimpleBankHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts)
// 	if err != nil {
// 		return err
// 	}

// 	err = userinfo.RegisterUserServiceHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts)
// 	if err != nil {
// 		return err
// 	}

// 	defer cancel()

// 	return http.Serve(listener, mux)

// }

var (
	grpcServerEndpoint = flag.String("grpc-server-endpoint", "localhost:9091", "gRPC server endpoint")
)

func main() {

	startGrpcServer()

}
