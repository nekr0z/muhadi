package db_test

import (
	"context"
	"crypto/sha256"
	"errors"
	"testing"

	"github.com/nekr0z/muhadi/internal/user"
	"github.com/stretchr/testify/assert"
)

var testUserName = "test_user"

func TestUser(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	testPasswordHash := sha256.Sum256([]byte("test_password"))

	t.Run("create user", func(t *testing.T) {
		err := testDB.NewUser(ctx, testUserName, testPasswordHash)
		assert.NoError(t, err)
	})

	t.Run("create conflicting user", func(t *testing.T) {
		anotherPasswordHash := sha256.Sum256([]byte("another_password"))
		err := testDB.NewUser(ctx, testUserName, anotherPasswordHash)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, user.ErrAlreadyExists))
	})

	t.Run("get hash", func(t *testing.T) {
		hash, err := testDB.GetHash(ctx, testUserName)
		assert.NoError(t, err)
		assert.Equal(t, testPasswordHash, hash)
	})
}
