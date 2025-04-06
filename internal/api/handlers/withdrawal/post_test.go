package withdrawal_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nekr0z/muhadi/internal/api/handlers/withdrawal"
	"github.com/nekr0z/muhadi/internal/api/handlers/withdrawal/mocks"
	"github.com/nekr0z/muhadi/internal/balance"
	"github.com/nekr0z/muhadi/internal/order"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap/zaptest"
)

var (
	orderID     = 2377225624
	sum         = 751.0
	requestBody = `{
"order": "2377225624",
"sum": 751
}`
)

func TestOrderPost(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ws := mocks.NewMockWithdrawalService(ctrl)

	log := zaptest.NewLogger(t)
	ts := mockTokenService{}

	post := withdrawal.WithdrawalHandleFunc(log, ws, ts)

	t.Run("unauthorized", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(requestBody))
		res := httptest.NewRecorder()

		post(res, req)

		assert.Equal(t, 401, res.Code)
		assert.True(t, ctrl.Satisfied())
	})

	t.Run("successful withdrawal", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(requestBody))
		req.AddCookie(cookie)
		res := httptest.NewRecorder()

		ws.EXPECT().Withdraw(gomock.Any(), testUserName, orderID, sum).Return(nil)

		post(res, req)

		assert.Equal(t, 200, res.Code)
		assert.True(t, ctrl.Satisfied())
	})

	t.Run("insufficient funds", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(requestBody))
		req.AddCookie(cookie)
		res := httptest.NewRecorder()

		ws.EXPECT().Withdraw(gomock.Any(), testUserName, orderID, sum).Return(balance.ErrNotEnoughFunds)

		post(res, req)

		assert.Equal(t, 402, res.Code)
		assert.True(t, ctrl.Satisfied())
	})

	t.Run("bad order ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(requestBody))
		req.AddCookie(cookie)
		res := httptest.NewRecorder()

		ws.EXPECT().Withdraw(gomock.Any(), testUserName, orderID, sum).Return(order.ErrOrderNumberInvalid)

		post(res, req)

		assert.Equal(t, 422, res.Code)
		assert.True(t, ctrl.Satisfied())
	})
}
