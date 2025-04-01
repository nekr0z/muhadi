package user_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/nekr0z/muhadi/internal/user"
	"github.com/stretchr/testify/assert"
)

func TestUserAuth(t *testing.T) {
	storage := make(mockStorage)
	svc := user.NewService(storage)

	t.Run("create", func(t *testing.T) {
		err := svc.NewUser(context.Background(), "test_user", "test_password")
		assert.NoError(t, err)
	})

	t.Run("auth", func(t *testing.T) {
		ok := svc.Auth(context.Background(), "test_user", "test_password")
		assert.True(t, ok)
	})

	t.Run("auth failed", func(t *testing.T) {
		ok := svc.Auth(context.Background(), "test_user", "another_password")
		assert.False(t, ok)
	})

	t.Run("unknown user", func(t *testing.T) {
		ok := svc.Auth(context.Background(), "unknown_user", "test_password")
		assert.False(t, ok)
	})

	t.Run("create existing user", func(t *testing.T) {
		err := svc.NewUser(context.Background(), "test_user", "another_password")
		assert.Error(t, err)
	})
}

type mockStorage map[string][32]byte

func (s mockStorage) NewUser(_ context.Context, userName string, hash [32]byte) error {
	if _, ok := s[userName]; ok {
		return user.ErrAlreadyExists
	}

	s[userName] = hash
	return nil
}

func (s mockStorage) GetHash(_ context.Context, userName string) ([32]byte, error) {
	v, ok := s[userName]
	if !ok {
		return [32]byte{}, fmt.Errorf("not found")
	}
	return v, nil
}
