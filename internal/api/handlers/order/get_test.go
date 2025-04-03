package order_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/nekr0z/muhadi/internal/api/handlers/order"
	"github.com/nekr0z/muhadi/internal/api/handlers/order/mocks"
	ordermodel "github.com/nekr0z/muhadi/internal/order"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap/zaptest"
)

func TestOrderGet(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	os := mocks.NewMockOrderService(ctrl)

	log := zaptest.NewLogger(t)
	ts := testTokenService{}
	cookie := &http.Cookie{Name: "token", Value: "test_token"}

	get := order.GetOrdersHandleFunc(log, os, ts)

	t.Run("empty", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(cookie)
		res := httptest.NewRecorder()

		os.EXPECT().GetOrders(gomock.Any(), testUserName).Return(nil, nil)

		get(res, req)

		assert.Equal(t, 204, res.Code)
		assert.True(t, ctrl.Satisfied())
	})

	t.Run("unauthorized", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()

		get(res, req)

		assert.Equal(t, 401, res.Code)
		assert.True(t, ctrl.Satisfied())
	})

	t.Run("orders", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(cookie)
		res := httptest.NewRecorder()

		t1, _ := time.Parse(time.RFC3339, "2020-12-10T15:15:45+03:00")
		t2, _ := time.Parse(time.RFC3339, "2020-12-10T15:12:01+03:00")
		t3, _ := time.Parse(time.RFC3339, "2020-12-09T16:09:53+03:00")
		testOrders := []*ordermodel.Order{
			{ID: 9278923470, Status: ordermodel.StatusDone, Accrual: 500, CreatedAt: t1},
			{ID: 12345678903, Status: ordermodel.StatusProcessing, CreatedAt: t2},
			{ID: 346436439, Status: ordermodel.StatusInvalid, CreatedAt: t3},
		}

		os.EXPECT().GetOrders(gomock.Any(), testUserName).Return(testOrders, nil)

		get(res, req)

		assert.Equal(t, 200, res.Code)
		assert.Equal(t, "application/json", res.Header().Get("Content-Type"))
		assert.JSONEq(t, `[
{
	"number": "9278923470",
	"status": "PROCESSED",
	"accrual": 500,
	"uploaded_at": "2020-12-10T15:15:45+03:00"
},
{
	"number": "12345678903",
	"status": "PROCESSING",
	"uploaded_at": "2020-12-10T15:12:01+03:00"
},
{
	"number": "346436439",
	"status": "INVALID",
	"uploaded_at": "2020-12-09T16:09:53+03:00"
}]`, res.Body.String())
		assert.True(t, ctrl.Satisfied())
	})
}
