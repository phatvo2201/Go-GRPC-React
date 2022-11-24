package implgrpc

import (
	"context"
	"fmt"
	"github.com/phatbb/wallet/config"
	"github.com/phatbb/wallet/models"
	wallet "github.com/phatbb/wallet/pb"
	"github.com/phatbb/wallet/service"
	"github.com/phatbb/wallet/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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
	wallet.UnimplementedAuthenServiceServer
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

func (as *AuthServer) SignUpUser(c context.Context, ui *wallet.SignUpUserRequest) (*wallet.SignUpUserResponse, error) {
	log.Println("Received request for adding repository with id " + fmt.Sprintf("%s", ui.Username))
	user := &models.SignUpInput{}

	user.Email = ui.Email
	user.Name = ui.Username
	user.Password = ui.Password
	user.PasswordConfirm = ui.PasswordConfirm

	// Logic to persist to database or storage.
	// working with mongodb
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	// Connect to MongoDB
	mongoconn := options.Client().ApplyURI(config.DBUri)
	// mongoconn := options.Client().ApplyURI(os.Getenv("MONGODB_CONNSTRING"))

	mongoclient, err := mongo.Connect(c, mongoconn)

	if err != nil {
		panic(err)
	}

	if err := mongoclient.Ping(c, readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("MongoDB successfully connected...")

	newUser, err := as.authService.SignUpUser(user)
	if err != nil {
		log.Println("cant not create usr at user implement ")

	}
	log.Println("Received request for adding repository with id " + fmt.Sprintf("%s %s has been created", newUser.Email, newUser.Name))

	//code := randstr.String(20)
	//verificationCode := utils.Encode(code)
	//updateData := &models.UpdateInput{
	//	VerificationCode: verificationCode,
	//}
	//as.userService.UpdateUser(newUser.ID.Hex(), updateData)
	////send email de user confirm
	//firstName := newUser.Name
	//
	//if strings.Contains(firstName, " ") {
	//	firstName = strings.Split(firstName, " ")[0]
	//}
	//
	////data for email service
	//
	//emailData := utils.EmailData{
	//	URL:       as.config.Origin + "/verifyemail/" + code,
	//	FirstName: firstName,
	//	Subject:   "verifycation code",
	//}
	//log.Println("this is phatbb1")
	//
	//err = utils.SendEmail(newUser, &emailData, "verificationCode.html")
	//if err != nil {
	//	return nil, status.Errorf(codes.Internal, "There was an error sending email: %s", err.Error())
	//
	//}

	newuser, err := as.authService.FindUserByEmail(user.Email)

	userId := newuser.ID
	walletuser := &models.CreateWalletRequest{}
	// oid, _ := primitive.ObjectIDFromHex(userId)
	walletuser.UserId = userId

	newWallet, err := as.authService.SignWallet(walletuser)
	if err != nil {
		log.Println("cant not creae wallet")
	}
	log.Printf("create wallet success", newWallet)

	return &wallet.SignUpUserResponse{
		User:  &wallet.User{Username: newUser.Name, Email: newUser.Email},
		Error: nil,
	}, nil

}

func (authServer *AuthServer) SignInUser(ctx context.Context, req *wallet.SignInUserRequest) (*wallet.SignInUserResponse, error) {
	user, err := authServer.authService.FindUserByEmail(req.GetEmail())
	log.Printf("this is your email %s", req.GetEmail())
	if err != nil {
		if err == mongo.ErrNoDocuments {

			return nil, status.Errorf(codes.InvalidArgument, "Invalid email ")

		}

		return nil, status.Errorf(codes.Internal, err.Error())

	}

	if err := utils.VerifyPassword(user.Password, req.GetPassword()); err != nil {

		return nil, status.Errorf(codes.InvalidArgument, "Invalid email or Password xxxx")

	}

	// Generate Tokens
	access_token, err := authServer.jwtManager.CreateToken(user)
	if err != nil {

		return nil, status.Errorf(codes.PermissionDenied, err.Error())

	}

	refresh_token, err := authServer.jwtManager.CreateReToken(user)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, err.Error())
	}

	res := &wallet.SignInUserResponse{
		Status:       "success",
		AccessToken:  access_token,
		RefreshToken: refresh_token,
	}

	return res, nil
}

func (as *AuthServer) FindUserById(ctx context.Context, in *wallet.GetInfoRequestId) (*wallet.UserResponse, error) {

	id := in.GetId()
	user, err := as.authService.FindUserById(id)
	if err != nil {
		return nil, err

	}
	return &wallet.UserResponse{
		User: &wallet.User{
			Username: user.Name,
			Email:    user.Email,
			Role:     user.Role,
			// Wallet: user.Wallet,

		},
	}, nil

}
func (as *AuthServer) RefreshToken(ctx context.Context, req *wallet.RefrehEmpty) (*wallet.SignInUserResponse, error) {
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

	claims, err := as.jwtManager.Verify(accessToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}
	log.Println(claims.Email)
	user, err := as.authService.FindUserByEmail(claims.Email)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "cannot find the email: %v", err)
	}
	access_token, err := as.jwtManager.CreateToken(user)
	if err != nil {

		return nil, status.Errorf(codes.PermissionDenied, err.Error())

	}

	res := &wallet.SignInUserResponse{
		Status:       "success",
		AccessToken:  access_token,
		RefreshToken: accessToken,
	}

	return res, nil
}
