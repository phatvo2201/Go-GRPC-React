package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/phatbb/auth/config"
	"github.com/phatbb/auth/implgrpc"
	"github.com/phatbb/auth/proto/auth"
	"github.com/phatbb/auth/service"
	"github.com/phatbb/auth/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

func startGrpcServer() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalln("cant not load config from env")
	}
	ctx := context.TODO()
	//create jwt manager
	jwtManger := utils.NewJwtManager(config.AccessTokenPublicKey, config.AccessTokenPrivateKey, 60*time.Minute)
	mongoConn := options.Client().ApplyURI(config.DBUri)
	mongoClient, err := mongo.Connect(ctx, mongoConn)
	if err != nil {
		log.Fatal("cant not connect to the mongodb database")
	}
	userCollection := mongoClient.Database(config.DBName).Collection("users")

	walletCollection := mongoClient.Database(config.DBName).Collection("wallet")
	authService := service.NewAuthService(userCollection, walletCollection, ctx)
	authServerHandler, _ := implgrpc.NewGrpcAuthServer(config, authService, userCollection, jwtManger)

	grpcServer := grpc.NewServer()
	auth.RegisterAuthenServiceServer(grpcServer, authServerHandler)

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
