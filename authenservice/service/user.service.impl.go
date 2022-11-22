package service

import (
	"context"
	"github.com/phatbb/wallet/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserServiceImpl struct {
	usercollection   *mongo.Collection
	walletcollection *mongo.Collection
	ctx              context.Context
}

func NewUserService(userollection *mongo.Collection, walletcollection *mongo.Collection, ctx context.Context) UserService {
	return &UserServiceImpl{userollection, walletcollection, ctx}
}

// SignInUsera implements UserService
func (uc *UserServiceImpl) FindUserByEmail(email string) (*models.DBResponse, error) {
	var user *models.DBResponse
	query := bson.M{"email": email}
	err := uc.usercollection.FindOne(uc.ctx, query).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &models.DBResponse{}, err
		}
		return nil, err
	}
	return user, nil

}

// SignUpUsera implements UserService
func (uc *UserServiceImpl) FindUserById(id string) (*models.DBResponse, error) {
	oid, _ := primitive.ObjectIDFromHex(id)

	var user *models.DBResponse
	query := bson.M{"_id": oid}
	err := uc.usercollection.FindOne(uc.ctx, query).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &models.DBResponse{}, err
		}
		return nil, err
	}
	return user, nil

}

// SignWalleta implements UserService
func (uc *UserServiceImpl) FindWalletByOwner(email string) (*models.DBWallet, error) {
	user, err := uc.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}
	wallet, err := uc.FindWallet(user.ID)
	if err != nil {
		return nil, err
	}
	return wallet, nil

}

// VerifyEmaila implements UserService

func (uc *UserServiceImpl) FindWallet(userId primitive.ObjectID) (*models.DBWallet, error) {
	var wallet *models.DBWallet

	query := bson.M{"user_id": userId}

	err := uc.walletcollection.FindOne(context.TODO(), query).Decode(&wallet)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			return nil, err
		}
		panic(err)
	}

	return wallet, nil

}
