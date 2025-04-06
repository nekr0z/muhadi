package order_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nekr0z/muhadi/internal/api/handlers/order"
	"github.com/nekr0z/muhadi/internal/api/handlers/order/mocks"
	ordermodel "github.com/nekr0z/muhadi/internal/order"
	orderservice "github.com/nekr0z/muhadi/internal/order/service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap/zaptest"
)

var (
	testUserName       = "test_user"
	correctOrderString = "12345678903"
	correctOrderID     = 12345678903
)

func TestOrderPost(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	os := mocks.NewMockOrderService(ctrl)

	log := zaptest.NewLogger(t)
	ts := testTokenService{}
	cookie := &http.Cookie{Name: "token", Value: "test_token"}

	post := order.PostOrderHandleFunc(log, os, ts)

	t.Run("create order", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(correctOrderString))
		req.AddCookie(cookie)
		res := httptest.NewRecorder()

		os.EXPECT().NewOrder(gomock.Any(), correctOrderID, testUserName).Return(nil)

		post(res, req)

		assert.Equal(t, 202, res.Code)
		assert.True(t, ctrl.Satisfied())
	})

	t.Run("resubmit order", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(correctOrderString))
		req.AddCookie(cookie)
		res := httptest.NewRecorder()

		os.EXPECT().NewOrder(gomock.Any(), correctOrderID, testUserName).Return(orderservice.ErrAlreadyExists)

		post(res, req)

		assert.Equal(t, 200, res.Code)
		assert.True(t, ctrl.Satisfied())
	})

	t.Run("invalid request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("ALL YOU BASE ARE BELONG TO US"))
		req.AddCookie(cookie)
		res := httptest.NewRecorder()

		post(res, req)

		assert.Equal(t, 400, res.Code)
		assert.True(t, ctrl.Satisfied())
	})

	t.Run("create order unauthorized", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("12345678903"))
		res := httptest.NewRecorder()
		post(res, req)

		assert.Equal(t, 401, res.Code)
	})

	t.Run("resubmit order submitted by another user", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(correctOrderString))
		req.AddCookie(cookie)
		res := httptest.NewRecorder()

		os.EXPECT().NewOrder(gomock.Any(), correctOrderID, testUserName).Return(orderservice.ErrOrderIDAlreadyTaken)

		post(res, req)

		assert.Equal(t, 409, res.Code)
		assert.True(t, ctrl.Satisfied())
	})

	t.Run("resubmit order submitted by another user", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("1234567890"))
		req.AddCookie(cookie)
		res := httptest.NewRecorder()

		os.EXPECT().NewOrder(gomock.Any(), 1234567890, testUserName).Return(ordermodel.ErrOrderNumberInvalid)

		post(res, req)

		assert.Equal(t, 422, res.Code)
		assert.True(t, ctrl.Satisfied())
	})
}

type testTokenService struct{}

func (t testTokenService) AuthUser(token string) (string, bool) {
	return testUserName, true
}
