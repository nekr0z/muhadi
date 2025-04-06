package balanceservice_test

import (
	"context"
	"errors"
	"testing"

	"github.com/nekr0z/muhadi/internal/balance"
	balanceservice "github.com/nekr0z/muhadi/internal/balance/service"
	"github.com/nekr0z/muhadi/internal/balance/service/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestBalanceService(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	st := mocks.NewMockStorage(ctrl)
	ag := &mockAccrualGetter{}

	svc := balanceservice.New(st, ag)

	t.Run("get balance", func(t *testing.T) {
		st.EXPECT().TotalWithdrawn(gomock.Any(), "test_user").Return(100.0, nil)
		b, w, err := svc.CurrentAndWithdrawnBalance(context.Background(), "test_user")
		assert.NoError(t, err)
		assert.Equal(t, 900.0, b)
		assert.Equal(t, 100.0, w)
		assert.True(t, ctrl.Satisfied())
	})

	t.Run(("successful withdrawal"), func(t *testing.T) {
		st.EXPECT().TotalWithdrawn(gomock.Any(), "test_user").Return(100.0, nil)
		st.EXPECT().SaveWithdrawal(gomock.Any(), gomock.Cond(func(w *balance.Withdrawal) bool {
			return w.UserName == "test_user" && w.OrderID == 12344 && w.Amount == 100.0
		})).Return(nil)
		err := svc.Withdraw(context.Background(), "test_user", 12344, 100.0)
		assert.NoError(t, err)
		assert.True(t, ctrl.Satisfied())
	})

	t.Run(("not enough funds"), func(t *testing.T) {
		st.EXPECT().TotalWithdrawn(gomock.Any(), "test_user").Return(100.0, nil)
		err := svc.Withdraw(context.Background(), "test_user", 12344, 2000.0)
		assert.True(t, errors.Is(err, balance.ErrNotEnoughFunds))
		assert.True(t, ctrl.Satisfied())
	})
}

type mockAccrualGetter struct {
}

func (m *mockAccrualGetter) TotalAccrual(_ context.Context, _ string) (float64, error) {
	return 1000, nil
}
