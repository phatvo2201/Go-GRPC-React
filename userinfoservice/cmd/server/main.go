package main

import (
	"context"
	"github.com/phatbb/userinfo/utils"
	"log"
	"net"

	"github.com/phatbb/userinfo/config"
	"github.com/phatbb/userinfo/implgrpc"
	"github.com/phatbb/userinfo/proto/userinfo"
	"github.com/phatbb/userinfo/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

func accessibleRole() map[string][]string {
	const infoservice = "/pb.UserService/"

	return map[string][]string{
		infoservice + "GetUserWalletInfo": {"user"},
		infoservice + "FindUserByEmail":   {"user"},
		infoservice + "FindUserById":      {"user"},
	}
}

func startGrpcServer() {
	log.Println("start auth grpc server")

	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalln("cant not load config from env")
	}
	ctx := context.TODO()
	//create jwt manager
	mongoConn := options.Client().ApplyURI(config.DBUri)
	mongoClient, err := mongo.Connect(ctx, mongoConn)
	if err != nil {
		log.Fatal("cant not connect to the mongodb database")
	}
	userCollection := mongoClient.Database(config.DBName).Collection("users")
	walletCollection := mongoClient.Database(config.DBName).Collection("wallet")

	userService := service.NewUserService(userCollection, walletCollection, ctx)

	userServerHandle := implgrpc.NewUserServerImpl(config, userService)

	jwtProvider := utils.NewJwtManager(config.AccessTokenPublicKey, config.AccessTokenPrivateKey, config.AccessTokenExpiresIn)

	//create and use interceptor
	interceptor := service.NewAuthInterceptor(jwtProvider, accessibleRole())

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.Unary()),
	)
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

func main() {

	startGrpcServer()

}
