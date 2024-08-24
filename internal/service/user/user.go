package user

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/mikesvis/gmart/internal/jwt"
	"github.com/mikesvis/gmart/internal/storage/user"
	"github.com/mikesvis/gmart/pkg/hash"
	"go.uber.org/zap"
)

var ErrBadRequest = errors.New(http.StatusText(http.StatusBadRequest))
var ErrUnauthorized = errors.New(http.StatusText(http.StatusUnauthorized))
var ErrConflict = errors.New(http.StatusText(http.StatusConflict))
var ErrIternal = errors.New(http.StatusText(http.StatusInternalServerError))

type Service struct {
	storage *user.Storage
	logger  *zap.SugaredLogger
}

func NewService(storage *user.Storage, logger *zap.SugaredLogger) *Service {
	return &Service{
		storage,
		logger,
	}
}

func (s *Service) RegisterUser(ctx context.Context, login, password string) (uint64, error) {
	if len(login) == 0 || len(password) == 0 {
		return 0, ErrBadRequest
	}

	s.logger.Infof("searching existing user by login %s", login)
	exist, err := s.storage.ExistByLogin(ctx, login)
	if err != nil {
		s.logger.Errorf("error while searching existing user by login %s, %v", login, err)
		return 0, ErrIternal
	}

	if exist {
		s.logger.Infof("user with login %s already exists", login)
		return 0, ErrConflict
	}

	s.logger.Infof("creating user with login %s", login)
	userID, err := s.storage.Create(ctx, login, hash.Hash([]byte(password)))
	if err != nil {
		s.logger.Errorf("error while creating new user by login %s, %v", login, err)
	}
	s.logger.Infof("created user with login %s id %d", login, userID)

	return userID, err
}

func (s *Service) Login(ctx context.Context, w http.ResponseWriter, userID uint64) error {

	expirationTime := time.Now().Add(jwt.TokenDuration)
	tokenString, err := jwt.CreateTokenString(userID, expirationTime)
	if err != nil {
		s.logger.Errorf("error while login user %d, %v", userID, err)
		return err
	}

	w.Header().Add(jwt.AuthorizationHeader, "Bearer "+tokenString)
	return nil
}

func (s *Service) GetUserID(ctx context.Context, login, password string) (uint64, error) {
	if len(login) == 0 || len(password) == 0 {
		return 0, ErrBadRequest
	}

	s.logger.Infof("searching user by creds (no password here he-he) login %s", login)
	userID, err := s.storage.GetUserID(ctx, login, hash.Hash([]byte(password)))

	if err != nil {
		s.logger.Errorf("error while searching user by creds with login %s, %v", login, err)
		return 0, ErrIternal
	}

	if userID == 0 {
		s.logger.Errorf("error because you are a cheater, %s!", login)
		return 0, ErrUnauthorized
	}

	s.logger.Infof("found user by creds login %s id %d", login, userID)
	return userID, nil
}
