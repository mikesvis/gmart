package middleware

import (
	_context "context"
	"errors"
	"net/http"
	"strings"

	"github.com/mikesvis/gmart/internal/context"
	"github.com/mikesvis/gmart/internal/jwt"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := getTokenFromHeader(r.Header.Get(jwt.AuthorizationHeader))
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		userID, err := jwt.GetUserIDFromTokenString(tokenString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := setUserIDToContext(r, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getTokenFromHeader(authHeader string) (string, error) {
	if len(authHeader) == 0 {
		return "", errors.New(http.StatusText(http.StatusUnauthorized))
	}

	parts := strings.Split(authHeader, " ")

	if len(parts) != 2 {
		return "", errors.New(http.StatusText(http.StatusUnauthorized))
	}

	tokenString := parts[1]

	if len(tokenString) == 0 {
		return "", errors.New(http.StatusText(http.StatusUnauthorized))
	}

	return tokenString, nil
}

func setUserIDToContext(r *http.Request, userID uint64) _context.Context {
	return _context.WithValue(r.Context(), context.UserIDContextKey, userID)
}
