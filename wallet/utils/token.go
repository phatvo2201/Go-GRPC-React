package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/phatbb/wallet/models"
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

func (jwtmanager *JWTManager) CreateToken(user *models.DBResponse) (string, error) {
	// decodedPrivateKey, err := base64.StdEncoding.DecodeString(jwtmanager.privateKey)
	// if err != nil {
	// 	return "", fmt.Errorf("could not decode key: %w", err)
	// }
	// key, err := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)

	// if err != nil {
	// 	return "", fmt.Errorf("create: parse key: %w", err)
	// }

	claims := UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(jwtmanager.tokenTime).Unix(),
		},
		Role:  user.Role,
		Email: user.Email,
	}

	// token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(jwtmanager.privateKey)

	// if err != nil {
	// 	return "", fmt.Errorf("create: sign token: %w", err)
	// }

	// return token, nil
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtmanager.privateKey))
}
func (jwtmanager *JWTManager) CreateReToken(user *models.DBResponse) (string, error) {
	// decodedPrivateKey, err := base64.StdEncoding.DecodeString(jwtmanager.privateKey)
	// if err != nil {
	// 	return "", fmt.Errorf("could not decode key: %w", err)
	// }
	// key, err := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)

	// if err != nil {
	// 	return "", fmt.Errorf("create: parse key: %w", err)
	// }

	claims := UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(jwtmanager.tokenTime * 1000).Unix(),
		},
		Role:  user.Role,
		Email: user.Email,
	}

	// token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(jwtmanager.privateKey)

	// if err != nil {
	// 	return "", fmt.Errorf("create: sign token: %w", err)
	// }

	// return token, nil
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtmanager.privateKey))
}
func (manager *JWTManager) Verify(accessToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected token signing method")
			}

			return []byte(manager.privateKey), nil
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

// func (jwtmanager JWTManager) ValidateToken(token string) (*UserClaims, error) {
// 	decodedPrivateKey, err := base64.StdEncoding.DecodeString(jwtmanager.privateKey)
// 	if err != nil {
// 		log.Panicln("1")

// 		return nil, fmt.Errorf("could not decode: %w", err)
// 		log.Panicln("1")
// 	}

// 	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)

// 	if err != nil {
// 		log.Panicln("2")

// 		return nil, fmt.Errorf("validate: parse key: %w", err)
// 	}

// 	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
// 		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
// 			log.Panicln("3")

// 			return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
// 		}
// 		return key, nil
// 	})

// 	if err != nil {
// 		return nil, fmt.Errorf("validate: %w", err)
// 	}

// 	claims, ok := parsedToken.Claims.(*UserClaims)
// 	if !ok || !parsedToken.Valid {
// 		log.Panicln("4")

// 		return nil, fmt.Errorf("validate: invalid token")
// 	}

// 	return claims, nil
// }
