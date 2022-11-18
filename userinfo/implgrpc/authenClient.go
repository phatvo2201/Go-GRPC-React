package implgrpc

import (
	"context"
	"log"
	"time"

	userinfo "github.com/phatbb/userinfo/pb"
	"google.golang.org/grpc"
)

type VerifyUserClient struct {
	service userinfo.SimpleBankClient
}

func NewVerifyUserClient(conn *grpc.ClientConn) *VerifyUserClient {
	service := userinfo.NewSimpleBankClient(conn)
	return &VerifyUserClient{service: service}
}

func (verifyClient *VerifyUserClient) VerifyOwnerByToken(in *userinfo.TokenAndEmail) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Millisecond*5000))
	defer cancel()

	res, err := verifyClient.service.VerifyOwner(ctx, in)

	if err != nil {
		log.Fatalf("GeMe: %v", err)
		return "cant not verify your token", err

	}

	return res.Status, nil
}
