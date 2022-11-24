package service

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"strings"
	"time"

	"github.com/phatbb/wallet/models"
	"github.com/phatbb/wallet/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AuthServiceImpl struct {
	usercollection   *mongo.Collection
	walletcollection *mongo.Collection
	ctx              context.Context
}

// SignWallet implements AuthService

func NewAuthService(userollection *mongo.Collection, walletcollection *mongo.Collection, ctx context.Context) *AuthServiceImpl {
	return &AuthServiceImpl{userollection, walletcollection, ctx}
}

func (uc *AuthServiceImpl) SignWallet(wallet *models.CreateWalletRequest) (*models.DBWallet, error) {
	wallet.Balance = 10000
	log.Println("inside phAT WALLET")
	wallet.CreateAt = time.Now()
	currency := "vietnamdong"
	wallet.Currency = fmt.Sprintf("%s", currency)
	wallet.UpdatedAt = wallet.CreateAt
	res, err := uc.walletcollection.InsertOne(uc.ctx, &wallet)
	if err != nil {
		if er, ok := err.(mongo.WriteException); ok && er.WriteErrors[0].Code == 11000 {
			return nil, errors.New("user with that email already exist")
		}
		return nil, err
	}
	var newWallet *models.DBWallet

	query := bson.M{"_id": res.InsertedID}

	err = uc.walletcollection.FindOne(uc.ctx, query).Decode(&newWallet)
	if err != nil {
		return nil, err
	}
	log.Println("da cai dat xong vi %s", newWallet.Currency)

	return newWallet, nil

}

func (as *AuthServiceImpl) SignUpUser(user *models.SignUpInput) (*models.DBResponse, error) {
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	user.Email = strings.ToLower(user.Email)
	user.Name = strings.ToLower(user.Name)

	user.PasswordConfirm = ""
	user.Verified = false
	user.Role = "user"

	hashedPassword, _ := utils.HashPassword(user.Password)
	user.Password = hashedPassword
	log.Println("before insert user to db")
	res, err := as.usercollection.InsertOne(as.ctx, &user)

	if err != nil {
		if er, ok := err.(mongo.WriteException); ok && er.WriteErrors[0].Code == 11000 {
			return nil, errors.New("user with that email already exist")
		}
		return nil, err
	}

	// Create a unique index for the email field
	opt := options.Index()
	opt.SetUnique(true)
	index := mongo.IndexModel{Keys: bson.M{"email": 1}, Options: opt}

	if _, err := as.usercollection.Indexes().CreateOne(as.ctx, index); err != nil {
		return nil, errors.New("could not create index for email")
	}

	var newUser *models.DBResponse
	query := bson.M{"_id": res.InsertedID}

	err = as.usercollection.FindOne(as.ctx, query).Decode(&newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}
func (as *AuthServiceImpl) SignInUser(user *models.SignUpInput) (*models.DBResponse, error) {
	return nil, nil

}
func (as *AuthServiceImpl) FindUserByEmail(email string) (*models.DBResponse, error) {
	var user *models.DBResponse
	query := bson.M{"email": email}
	err := as.usercollection.FindOne(as.ctx, query).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &models.DBResponse{}, err
		}
		return nil, err
	}
	return user, nil

}
func (as *AuthServiceImpl) FindUserById(id string) (*models.DBResponse, error) {
	oid, _ := primitive.ObjectIDFromHex(id)

	var user *models.DBResponse
	query := bson.M{"_id": oid}
	err := as.usercollection.FindOne(as.ctx, query).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &models.DBResponse{}, err
		}
		return nil, err
	}
	return user, nil

}
