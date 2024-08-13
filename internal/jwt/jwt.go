package jwt

import (
	"errors"
	"time"

	_jwt "github.com/golang-jwt/jwt/v5"
)

const SecretPass = "mySecretPass"
const AuthorizationHeader = "Authorization"
const TokenDuration = time.Hour * 24

type Claims struct {
	UserID uint64 `json:"userId"`
	_jwt.RegisteredClaims
}

func CreateTokenString(userID uint64, exp time.Time) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: _jwt.RegisteredClaims{
			ExpiresAt: _jwt.NewNumericDate(exp),
		},
	}
	token := _jwt.NewWithClaims(_jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(SecretPass))
	if err != nil {
		return "", err
	}

	return tokenString, err
}

func GetUserIDFromTokenString(tokenString string) (uint64, error) {
	claims := &Claims{}

	token, err := _jwt.ParseWithClaims(tokenString, claims, func(token *_jwt.Token) (any, error) {
		return []byte(SecretPass), nil
	})

	if err == nil && token.Valid {
		if claims.UserID == 0 {
			return 0, errors.New(`undefined user`)
		}

		return claims.UserID, nil
	}

	if err == nil && !token.Valid {
		return 0, errors.New(`invalid token`)
	}

	return 0, err
}
