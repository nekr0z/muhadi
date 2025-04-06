package balance_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nekr0z/muhadi/internal/api/handlers/balance"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestOrderGet(t *testing.T) {
	t.Parallel()

	log := zaptest.NewLogger(t)
	ts := mockTokenService{}
	bs := mockBalanceService{}
	cookie := &http.Cookie{Name: "token", Value: "test_token"}

	get := balance.BalanceHandleFunc(log, bs, ts)

	t.Run("unauthorized", func(t *testing.T) {
		assert.HTTPStatusCode(t, get, http.MethodGet, "/", nil, 401)
	})

	t.Run("balance", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(cookie)
		res := httptest.NewRecorder()

		get(res, req)

		assert.Equal(t, 200, res.Code)
		assert.Equal(t, "application/json", res.Header().Get("Content-Type"))
		assert.JSONEq(t, `{
"current": 500.5,
"withdrawn": 42
    }`, res.Body.String())
	})
}

var testUserName = "test_user"

type mockTokenService struct{}

func (t mockTokenService) AuthUser(token string) (string, bool) {
	return testUserName, true
}

type mockBalanceService struct{}

func (b mockBalanceService) CurrentAndWithdrawnBalance(_ context.Context, userName string) (float64, float64, error) {
	return 500.5, 42, nil
}
