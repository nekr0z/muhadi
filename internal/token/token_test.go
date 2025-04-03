package token_test

import (
	"testing"

	"github.com/nekr0z/muhadi/internal/token"
	"github.com/stretchr/testify/assert"
)

func TestToken(t *testing.T) {
	t.Parallel()

	ts := token.New()

	userName := "test_user"

	token, err := ts.NewToken(userName)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	name, ok := ts.AuthUser(token)
	assert.True(t, ok)
	assert.Equal(t, userName, name)
}
