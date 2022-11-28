package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/phatbb/auth/models"
)

type JWTManager struct {
	PublicKey  string
	privateKey string
	tokenTime  time.Duration
}

type UserClaims struct {
	jwt.StandardClaims
	Email string `json:"email"`
	Role  string `json:"role"`
}

func NewJwtManager(pbk string, privatekey string, duration time.Duration) *JWTManager {
	return &JWTManager{pbk, privatekey, duration}
}

func (jm *JWTManager) CreateToken(user *models.DBResponse) (string, error) {

	claims := UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(jm.tokenTime).Unix(),
		},
		Role:  user.Role,
		Email: user.Email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jm.privateKey))
}
func (jm *JWTManager) CreateReToken(user *models.DBResponse) (string, error) {

	claims := UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(jm.tokenTime * 1000).Unix(),
		},
		Role:  user.Role,
		Email: user.Email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jm.privateKey))
}
func (jm *JWTManager) Verify(accessToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected token signing method")
			}

			return []byte(jm.privateKey), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
