package service

import (
	"github.com/phatbb/wallet/models"
)

type UserService interface {
	FindUserByEmail(email string) (*models.DBResponse, error)
	FindUserById(id string) (*models.DBResponse, error)
	FindWalletByOwner(gmail string) (*models.DBWallet, error)
}
