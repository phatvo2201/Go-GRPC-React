package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateWalletRequest struct {
	Balance   int                `json:"balance" bson:"balance" binding:"required"`
	Currency  string             `json:"currency" bson:"currency"  binding:"required"`
	UserId    primitive.ObjectID `json:"user_id,omitempty" bson:"user_id"  binding:"required"`
	CreateAt  time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type DBWallet struct {
	Id        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Balance   int                `json:"balance" bson:"balance" binding:"required"`
	Currency  string             `json:"currency" bson:"currency" binding:"required"`
	User      string             `json:"user" bson:"user" binding:"required"`
	CreateAt  time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
