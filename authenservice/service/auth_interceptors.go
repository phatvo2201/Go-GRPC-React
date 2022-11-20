package service

import (
	"context"
	"fmt"
	"log"
	"strings"

	wallet "github.com/phatbb/wallet/pb"
	"github.com/phatbb/wallet/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

//erver interceptor for authentication and authorization
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

		// type Req struct {
		// 	email string
		// 	role  string
		// }
		xType := fmt.Sprintf("%T aaaaaaaaaaaaaaaa", req)
		log.Println(xType) // "[]int"
		_, ok := req.(*wallet.GetInfoRequestGmail)
		if ok {
			log.Println("intercept for get email owner")
			log.Println("--> unary interceptor: ", req.(*wallet.GetInfoRequestGmail).Gmail)
			err := interceptor.authorize(ctx, info.FullMethod, req.(*wallet.GetInfoRequestGmail).Gmail)
			if err != nil {
				return nil, err
			}

		}

		log.Println("--> unary interceptor: ", info.FullMethod)

		// Use the provided zap.Logger for logging but use the fields from context.

		err := interceptor.authorize(ctx, info.FullMethod, "")
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}
func (interceptor *AuthInterceptor) authorize(ctx context.Context, method string, email string) error {
	log.Println("sadsadsadasdasdasd")
	accessibleRoles, ok := interceptor.accessibleRoles[method]
	if !ok {
		// everyone can access
		log.Println("everyone can access")
		return nil
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Println("this is all can access the grpc ")
		return status.Errorf(codes.Unauthenticated, "access token is not provider1")

	}
	log.Println("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx ")
	autheninfo := md["authorization"]
	// log.Printf("this is the metadata %s", md)
	if len(autheninfo) == 0 {
		log.Println("this is all can access the grpc ")
		return status.Errorf(codes.Unauthenticated, "access token is not provider2 %s xzczxczxc", method)
	}

	log.Println("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx ")

	accessToken := autheninfo[0]

	// using the function
	if strings.Contains(accessToken, "Bearer ") {
		accessToken = accessToken[7:]
	}

	log.Printf("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx %s xxxxxxxxccasda ", accessToken)

	claims, err := interceptor.jwtManager.Verify(accessToken)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}
	log.Println(claims.Email)

	log.Printf("hellooo %s", email)

	log.Printf("hellooo %s", claims.Email)

	if email != "" {
		if claims.Email != email {
			log.Println("helloooo")
			return status.Errorf(codes.Unauthenticated, "you can not have permissccess this info : %v", err)

		}

	}
	for _, role := range accessibleRoles {
		if role == claims.Role {
			return nil
		}
	}

	return status.Error(codes.PermissionDenied, "no permission to access this RPC")
}
