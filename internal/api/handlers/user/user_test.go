package user_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nekr0z/muhadi/internal/api/handlers"
	"github.com/nekr0z/muhadi/internal/api/handlers/user"
	"github.com/nekr0z/muhadi/internal/token"
	userservice "github.com/nekr0z/muhadi/internal/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestUser(t *testing.T) {
	t.Parallel()

	log := zaptest.NewLogger(t)
	ts := token.New()
	us := make(testUserService)

	reg := user.RegisterHandleFunc(log, us, ts)
	login := user.LoginHandleFunc(log, us, ts)

	var token string
	t.Run("create user", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"login": "test_user", "password": "test_password"}`))
		res := httptest.NewRecorder()
		reg(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		response := res.Result()
		defer response.Body.Close()

		c := response.Cookies()

		require.Len(t, c, 1)
		require.Equal(t, handlers.TokenCookieName, c[0].Name)
		token = c[0].Value
	})

	t.Run("create user with same name", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"login": "test_user", "password": "password"}`))
		res := httptest.NewRecorder()
		reg(res, req)

		assert.Equal(t, http.StatusConflict, res.Code)
	})

	t.Run("login", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"login": "test_user", "password": "test_password"}`))
		res := httptest.NewRecorder()
		login(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		response := res.Result()
		defer response.Body.Close()

		c := response.Cookies()
		require.Len(t, c, 1)
		require.Equal(t, handlers.TokenCookieName, c[0].Name)
		require.Equal(t, token, c[0].Value)
	})

	t.Run("login with wrong password", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"login": "test_user", "password": "password"}`))
		res := httptest.NewRecorder()
		login(res, req)

		assert.Equal(t, http.StatusUnauthorized, res.Code)

		response := res.Result()
		defer response.Body.Close()

		c := response.Cookies()
		require.Len(t, c, 0)
	})
}

type testUserService map[string]string

func (t testUserService) NewUser(_ context.Context, userName, password string) error {
	if _, ok := t[userName]; ok {
		return userservice.ErrAlreadyExists
	}

	t[userName] = password
	return nil
}

func (t testUserService) Auth(_ context.Context, userName, password string) bool {
	return t[userName] == password
}
