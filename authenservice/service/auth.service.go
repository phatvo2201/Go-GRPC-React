package service

import (
	"github.com/phatbb/wallet/models"
)

type AuthService interface {
	SignUpUser(*models.SignUpInput) (*models.DBResponse, error)
	SignInUser(*models.SignUpInput) (*models.DBResponse, error)

	SignWallet(user *models.CreateWalletRequest) (*models.DBWallet, error)
	VerifyEmail(code string) (*models.DBResponse, error)
}
