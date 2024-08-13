package accural

import "github.com/mikesvis/gmart/internal/storage/accural"

type Service struct {
	storage *accural.Storage
}

func NewService(storage *accural.Storage) *Service {
	return &Service{
		storage,
	}
}
