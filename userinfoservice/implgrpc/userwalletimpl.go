package implgrpc

import (
	"context"
	"log"
	"strings"

	"github.com/phatbb/userinfo/config"
	userinfo "github.com/phatbb/userinfo/pb"
	"github.com/phatbb/userinfo/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type UserServer struct {
	userService service.UserService
	config      config.Config
	grpcClient  *VerifyUserClient
	userinfo.UnimplementedUserServiceServer
}

func NewUserServerImpl(config config.Config, userService service.UserService, grpcClient *VerifyUserClient) *UserServer {
	return &UserServer{
		userService: userService,
		config:      config,
		grpcClient:  grpcClient,
	}
}

func (us *UserServer) FindUserById(ctx context.Context, in *userinfo.GetInfoRequestId) (*userinfo.User, error) {

	id := in.GetId()
	user, err := us.userService.FindUserById(id)
	if err != nil {
		return nil, err

	}
	return &userinfo.User{

		Username: user.Name,
		Email:    user.Email,
		Role:     user.Role,
		// Wallet: user.Wallet,

	}, nil

}
func (us *UserServer) FindUserByEmail(ctx context.Context, in *userinfo.GetInfoRequestGmail) (*userinfo.User, error) {
	email := in.GetGmail()
	user, err := us.userService.FindUserByEmail(email)
	if err != nil {
		return nil, err

	}
	return &userinfo.User{

		Username: user.Name,
		Email:    user.Email,
		Role:     user.Role,
		// Wallet: user.Wallet,

	}, nil

}
func (us *UserServer) GetUserWalletInfo(ctx context.Context, in *userinfo.GetInfoRequestGmail) (*userinfo.Wallet, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	log.Printf("this is the metadataaaaa %s", md)
	if !ok {
		log.Println("this is all can access the grpc ")
		return nil, status.Errorf(codes.Unauthenticated, "access token is not provider1")

	}
	autheninfo := md["authorization"]
	if len(autheninfo) == 0 {
		log.Println("this is all can access the grpc ")
		return nil, status.Errorf(codes.Unauthenticated, "access token is not provider2222")
	}
	accessToken := autheninfo[0]

	// using the function
	if strings.Contains(accessToken, "Bearer ") {
		accessToken = accessToken[7:]
	}

	log.Printf("yyyyyyyyyyyyyyyyyyy %s yyyyy ", accessToken)
	email := in.GetGmail()
	log.Printf("this is the email of user %s", email)

	req := &userinfo.TokenAndEmail{
		Token: accessToken,
		Email: email,
	}
	status, err := us.grpcClient.VerifyOwnerByToken(req)
	if err != nil {
		return nil, err
	}
	log.Printf("verify success %s", status)

	walletuser, err := us.userService.FindWalletByOwner(email)
	if err != nil {
		log.Println("error when get user by email")
		return nil, err
	}
	return &userinfo.Wallet{
		Balance:  uint64(walletuser.Balance),
		Currency: walletuser.Currency,
	}, nil
}
