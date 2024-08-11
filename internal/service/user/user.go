package user

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/mikesvis/gmart/internal/jwt"
	"github.com/mikesvis/gmart/internal/storage/user"
	"github.com/mikesvis/gmart/pkg/hash"
)

var ErrBadRequest = errors.New(http.StatusText(http.StatusBadRequest))
var ErrUnauthorized = errors.New(http.StatusText(http.StatusUnauthorized))
var ErrConflict = errors.New(http.StatusText(http.StatusConflict))
var ErrIternal = errors.New(http.StatusText(http.StatusInternalServerError))

type Service struct {
	storage *user.Storage
}

func NewService(storage *user.Storage) *Service {
	return &Service{
		storage,
	}
}

func (s *Service) RegisterUser(ctx context.Context, login, password string) (uint64, error) {
	if len(login) == 0 || len(password) == 0 {
		return 0, ErrBadRequest
	}

	exist, err := s.storage.ExistByLogin(ctx, login)
	if err != nil {
		return 0, ErrIternal
	}

	if exist {
		return 0, ErrConflict
	}

	return s.storage.Create(ctx, login, hash.Hash([]byte(password)))
}

func (s *Service) Login(ctx context.Context, w http.ResponseWriter, userID uint64) error {

	expirationTime := time.Now().Add(jwt.TokenDuration)
	tokenString, err := jwt.CreateTokenString(userID, expirationTime)
	if err != nil {
		return err
	}

	w.Header().Add("Authorization", "Bearer "+tokenString)
	return nil
}

func (s *Service) GetUserID(ctx context.Context, login, password string) (uint64, error) {
	if len(login) == 0 || len(password) == 0 {
		return 0, ErrBadRequest
	}

	userID, err := s.storage.GetUserID(ctx, login, hash.Hash([]byte(password)))

	if err != nil {
		return 0, ErrIternal
	}

	if userID == 0 {
		return 0, ErrUnauthorized
	}

	return userID, nil
}
