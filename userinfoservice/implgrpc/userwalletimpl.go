package implgrpc

import (
	"context"
	"github.com/phatbb/userinfo/config"
	"github.com/phatbb/userinfo/proto/userinfo"
	"github.com/phatbb/userinfo/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"strconv"
)

type UserServer struct {
	userService *service.UserServiceImpl
	config      config.Config
	userinfo.UnimplementedUserServiceServer
}

func NewUserServerImpl(config config.Config, userService *service.UserServiceImpl) *UserServer {
	return &UserServer{
		userService: userService,
		config:      config,
	}
}

func (us *UserServer) FindUserById(ctx context.Context, in *userinfo.GetInfoRequestId) (*userinfo.User, error) {

	id := in.GetId()
	user, err := us.userService.FindUserById(id)
	if err != nil {
		log.Println("can not find user by id")
		return nil, err

	}
	return &userinfo.User{

		Username: user.Name,
		Email:    user.Email,
		Role:     user.Role,
	}, nil

}
func (us *UserServer) FindUserByEmail(ctx context.Context, in *userinfo.GetInfoRequestGmail) (*userinfo.User, error) {
	email := in.GetGmail()
	user, err := us.userService.FindUserByEmail(email)
	if err != nil {
		log.Println("can not find user by email")
		return nil, err

	}
	return &userinfo.User{

		Username: user.Name,
		Email:    user.Email,
		Role:     user.Role,
	}, nil

}
func (us *UserServer) GetUserWalletInfo(ctx context.Context, in *userinfo.GetInfoRequestGmail) (*userinfo.Wallet, error) {
	email := in.GetGmail()
	//header := metadata.New(map[string]string{"Cookie": cookie})
	header := metadata.Pairs("gauserId", strconv.Itoa(1111111111111111))
	if err := grpc.SendHeader(ctx, header); err != nil {
		return nil, status.Errorf(codes.Internal, "unable to send 'x-response-id' header")
	}

	walletUser, err := us.userService.FindWalletByOwner(email)
	if err != nil {
		log.Println("error when get user by email")
		return nil, err
	}
	return &userinfo.Wallet{
		Balance:  uint64(walletUser.Balance),
		Currency: walletUser.Currency,
	}, nil
}
