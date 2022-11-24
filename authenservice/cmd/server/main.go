package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/phatbb/wallet/config"
	"github.com/phatbb/wallet/implgrpc"
	wallet "github.com/phatbb/wallet/pb"
	"github.com/phatbb/wallet/service"
	"github.com/phatbb/wallet/utils"
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
	jwtManger := utils.NewJwtManager(config.AccessTokenPublicKey, config.AccessTokenPrivateKey, 1500*time.Minute)
	mongoconn := options.Client().ApplyURI(config.DBUri)
	mongoClient, err := mongo.Connect(ctx, mongoconn)
	if err != nil {
		log.Fatal("cant not connect to the mongodb database")
	}
	usercollection := mongoClient.Database(config.DBName).Collection("users")

	walletCollection := mongoClient.Database(config.DBName).Collection("wallet")
	authService := service.NewAuthService(usercollection, walletCollection, ctx)
	authServerHandler, _ := implgrpc.NewGrpcAuthServer(config, authService, usercollection, jwtManger)

	grpcServer := grpc.NewServer()
	wallet.RegisterAuthenServiceServer(grpcServer, authServerHandler)

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
