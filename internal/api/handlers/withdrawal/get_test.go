package withdrawal_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/nekr0z/muhadi/internal/api/handlers/withdrawal"
	"github.com/nekr0z/muhadi/internal/api/handlers/withdrawal/mocks"
	"github.com/nekr0z/muhadi/internal/balance"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap/zaptest"
)

func TestGetWithdrawals(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	os := mocks.NewMockWithdrawalService(ctrl)

	log := zaptest.NewLogger(t)
	ts := mockTokenService{}

	get := withdrawal.GetWithdrawalsHandleFunc(log, os, ts)

	t.Run("empty", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(cookie)
		res := httptest.NewRecorder()

		os.EXPECT().GetWithdrawals(gomock.Any(), testUserName).Return(nil, nil)

		get(res, req)

		assert.Equal(t, 204, res.Code)
		assert.True(t, ctrl.Satisfied())
	})

	t.Run("unauthorized", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()

		get(res, req)

		assert.Equal(t, 401, res.Code)
	})

	t.Run("orders", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(cookie)
		res := httptest.NewRecorder()

		t1, _ := time.Parse(time.RFC3339, "2020-12-09T16:09:57+03:00")
		testWithdrawals := []balance.Withdrawal{
			{UserName: testUserName, OrderID: 2377225624, Amount: 500, At: t1},
		}

		os.EXPECT().GetWithdrawals(gomock.Any(), testUserName).Return(testWithdrawals, nil)

		get(res, req)

		assert.Equal(t, 200, res.Code)
		assert.Equal(t, "application/json", res.Header().Get("Content-Type"))
		assert.JSONEq(t, `[        {
"order": "2377225624",
"sum": 500,
"processed_at": "2020-12-09T16:09:57+03:00"
}]`, res.Body.String())
		assert.True(t, ctrl.Satisfied())
	})
}
