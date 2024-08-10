package jwt

import (
	"time"

	_jwt "github.com/golang-jwt/jwt/v5"
)

const SecretPass = "mySecretPass"
const AuthorizationCookieName = "Authorization-JWT"
const TokenDuration = time.Hour * 24

type Claims struct {
	UserID uint64 `json:"userId"`
	_jwt.RegisteredClaims
}

// func GetUserIDFromTokenString(tokenString string) (string, error) {
// 	claims := &Claims{}

// 	token, err := _jwt.ParseWithClaims(tokenString, claims, func(token *_jwt.Token) (any, error) {
// 		return []byte(SecretPass), nil
// 	})

// 	// все хорошо в куке, не трогаем
// 	if err == nil && token.Valid {
// 		// пустой UserID в токене (по заданию)
// 		if claims.UserID == 0 {
// 			return "", errors.ErrEmptyUserID
// 		}

// 		return claims.UserID, nil
// 	}

// 	if err == nil && !token.Valid {
// 		return "", errors.ErrInvalidToken
// 	}

// 	return "", err
// }

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
