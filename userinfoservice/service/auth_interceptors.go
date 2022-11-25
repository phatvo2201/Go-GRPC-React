package service

import (
	"context"
	"log"
	"strings"

	"github.com/phatbb/userinfo/proto/userinfo"
	"github.com/phatbb/userinfo/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	jwtManager      *utils.JWTManager
	accessibleRoles map[string][]string
}

func NewAuthInterceptor(jwtManager *utils.JWTManager, accessibleRoles map[string][]string) *AuthInterceptor {
	return &AuthInterceptor{jwtManager, accessibleRoles}
}

// Unary returns a server interceptor function to authenticate and authorize unary RPC
func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		_, ok := req.(*userinfo.GetInfoRequestGmail)
		if ok {

			log.Println("--> unary interceptor: ", req.(*userinfo.GetInfoRequestGmail).Gmail)
			err := interceptor.authorize(ctx, info.FullMethod, req.(*userinfo.GetInfoRequestGmail).Gmail)
			if err != nil {
				return nil, err
			}

		}

		err := interceptor.authorize(ctx, info.FullMethod, "")
		if err != nil {
			log.Printf("you have an error protect by interceptor ,this err is %s", err)
			return nil, err
		}

		return handler(ctx, req)
	}
}
func (interceptor *AuthInterceptor) authorize(ctx context.Context, method string, email string) error {
	accessibleRoles, ok := interceptor.accessibleRoles[method]
	if !ok {
		// everyone can access
		log.Println("everyone can access")
		return nil
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "can not access context")

	}
	autheninfo := md["authorization"]
	if len(autheninfo) == 0 {
		return status.Errorf(codes.Unauthenticated, "access token is not provider %s", method)
	}

	accessToken := autheninfo[0]

	// using the function
	if strings.Contains(accessToken, "Bearer ") {
		accessToken = accessToken[7:]
	}

	claims, err := interceptor.jwtManager.Verify(accessToken)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "access token is invalid by validate %s", err)
	}

	if email != "" {
		if claims.Email != email {
			log.Println("the email not right")
			return status.Errorf(codes.Unauthenticated, "you can not have permission this info : %v", err)

		}

	}
	//check role in token and in acessrole
	for _, role := range accessibleRoles {
		if role == claims.Role {
			return nil
		}
	}

	return status.Error(codes.PermissionDenied, "no permission to access this RPC")
}
