package orderservice_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/nekr0z/muhadi/internal/order"
	orderservice "github.com/nekr0z/muhadi/internal/order/service"
	"github.com/nekr0z/muhadi/internal/order/service/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestOrderService(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	st := mocks.NewMockStorage(ctrl)
	svc := orderservice.New(st)

	t.Run("new order", func(t *testing.T) {
		st.EXPECT().CreateOrder(gomock.Any(), gomock.Cond(func(o *order.Order) bool {
			return o.ID == 12344 && o.Status == order.StatusNew && o.UserName == "test_user"
		})).Return(nil)

		err := svc.NewOrder(context.Background(), 12344, "test_user")
		assert.NoError(t, err)
		assert.True(t, ctrl.Satisfied())
	})

	t.Run("existing order", func(t *testing.T) {
		st.EXPECT().CreateOrder(gomock.Any(), gomock.Cond(func(o *order.Order) bool {
			return o.ID == 12344 && o.Status == order.StatusNew && o.UserName == "test_user"
		})).Return(fmt.Errorf("already exists"))
		st.EXPECT().GetOrder(gomock.Any(), 12344).Return(order.New(12344, "test_user"), nil)

		err := svc.NewOrder(context.Background(), 12344, "test_user")
		assert.True(t, errors.Is(err, orderservice.ErrAlreadyExists))
		assert.True(t, ctrl.Satisfied())
	})

	t.Run("order submitted by another user", func(t *testing.T) {
		st.EXPECT().CreateOrder(gomock.Any(), gomock.Cond(func(o *order.Order) bool {
			return o.ID == 12344 && o.Status == order.StatusNew && o.UserName == "another_user"
		})).Return(fmt.Errorf("already exists"))
		st.EXPECT().GetOrder(gomock.Any(), 12344).Return(order.New(12344, "test_user"), nil)

		err := svc.NewOrder(context.Background(), 12344, "another_user")
		assert.True(t, errors.Is(err, orderservice.ErrOrderIDAlreadyTaken))
		assert.True(t, ctrl.Satisfied())
	})

	t.Run("invalid order", func(t *testing.T) {
		err := svc.NewOrder(context.Background(), 12345, "test_user")
		assert.True(t, errors.Is(err, order.ErrOrderNumberInvalid))
		assert.True(t, ctrl.Satisfied())
	})
}
