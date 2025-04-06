package user

import (
	"context"
	"crypto/sha256"
	"errors"
)

type Service struct {
	storage Storage
}

func NewService(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) NewUser(ctx context.Context, userName, password string) error {
	hash := sha256.Sum256([]byte(password))

	return s.storage.NewUser(ctx, userName, hash)
}

func (s *Service) Auth(ctx context.Context, userName, password string) bool {
	hash := sha256.Sum256([]byte(password))

	storedHash, err := s.storage.GetHash(ctx, userName)
	if err != nil {
		return false
	}

	return hash == storedHash
}

type Storage interface {
	NewUser(ctx context.Context, userName string, hash [32]byte) error
	GetHash(ctx context.Context, userName string) ([32]byte, error)
}

var ErrAlreadyExists = errors.New("user already exists")
