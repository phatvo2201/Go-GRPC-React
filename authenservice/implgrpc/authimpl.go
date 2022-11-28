package implgrpc

import (
	"context"
	"github.com/phatbb/auth/config"
	"github.com/phatbb/auth/models"
	"github.com/phatbb/auth/proto/auth"
	"github.com/phatbb/auth/proto/userinfo"

	"github.com/phatbb/auth/service"
	"github.com/phatbb/auth/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"strings"
)

type AuthServer struct {
	config         config.Config
	authService    *service.AuthServiceImpl
	userCollection *mongo.Collection
	jwtManager     *utils.JWTManager
	auth.UnimplementedAuthenServiceServer
}

func NewGrpcAuthServer(config config.Config, authService *service.AuthServiceImpl,
	userCollection *mongo.Collection, jwtManager *utils.JWTManager) (*AuthServer, error) {

	authServer := &AuthServer{
		config:         config,
		authService:    authService,
		userCollection: userCollection,
		jwtManager:     jwtManager,
	}

	return authServer, nil
}

func (as *AuthServer) SignUpUser(c context.Context, ui *auth.SignUpUserRequest) (*auth.SignUpUserResponse, error) {
	user := &models.SignUpInput{}

	user.Email = ui.Email
	user.Name = ui.Username
	user.Password = ui.Password
	user.PasswordConfirm = ui.PasswordConfirm

	newUser, err := as.authService.SignUpUser(user)
	if err != nil {
		log.Println("cant not create usr at user implement ")
		return nil, err

	}

	userId := newUser.ID
	walletUser := &models.CreateWalletRequest{}
	walletUser.UserId = userId

	newWallet, err := as.authService.SignWallet(walletUser)
	if err != nil {
		log.Println("cant not creae wallet")
		return nil, err
	}
	log.Printf("create wallet success", newWallet)

	return &auth.SignUpUserResponse{
		User:  &userinfo.User{Username: newUser.Name, Email: newUser.Email},
		Error: nil,
	}, nil

}

func (as *AuthServer) SignInUser(ctx context.Context, req *auth.SignInUserRequest) (*auth.SignInUserResponse, error) {
	user, err := as.authService.FindUserByEmail(req.GetEmail())
	if err != nil {
		if err == mongo.ErrNoDocuments {

			return nil, status.Errorf(codes.InvalidArgument, "Invalid email ")

		}
		return nil, status.Errorf(codes.Internal, err.Error())

	}

	if err := utils.VerifyPassword(user.Password, req.GetPassword()); err != nil {

		return nil, status.Errorf(codes.InvalidArgument, "Invalid email or Password")

	}

	// Generate Tokens
	accessToken, err := as.jwtManager.CreateToken(user)
	if err != nil {

		return nil, status.Errorf(codes.PermissionDenied, err.Error())

	}

	refreshToken, err := as.jwtManager.CreateReToken(user)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, err.Error())
	}

	res := &auth.SignInUserResponse{
		Status:       "success",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return res, nil
}

func (as *AuthServer) FindUserById(ctx context.Context, in *userinfo.GetInfoRequestId) (*userinfo.UserResponse, error) {

	id := in.GetId()
	user, err := as.authService.FindUserById(id)
	if err != nil {
		return nil, err

	}
	return &userinfo.UserResponse{
		User: &userinfo.User{
			Username: user.Name,
			Email:    user.Email,
			Role:     user.Role,
		},
	}, nil

}
func (as *AuthServer) RefreshToken(ctx context.Context, req *auth.RefrehEmpty) (*auth.SignInUserResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Println("this is all can access the grpc ")
		return nil, status.Errorf(codes.Unauthenticated, "access metadata fail")

	}
	authInfo := md["authorization"]
	if len(authInfo) == 0 {
		log.Println("this is all can access the grpc ")
		return nil, status.Errorf(codes.Unauthenticated, "access token is not provider")
	}
	accessToken := authInfo[0]

	// using the function
	if strings.Contains(accessToken, "Bearer ") {
		accessToken = accessToken[7:]
	}

	claims, err := as.jwtManager.Verify(accessToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}
	user, err := as.authService.FindUserByEmail(claims.Email)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "cannot find the email: %v", err)
	}
	newAccessToken, err := as.jwtManager.CreateToken(user)
	if err != nil {

		return nil, status.Errorf(codes.PermissionDenied, err.Error())

	}
	newRefreshToken, err := as.jwtManager.CreateReToken(user)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, err.Error())
	}

	res := &auth.SignInUserResponse{
		Status:       "success",
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}

	return res, nil
}
