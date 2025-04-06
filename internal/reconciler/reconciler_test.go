package reconciler_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/nekr0z/muhadi/internal/order"
	"github.com/nekr0z/muhadi/internal/reconciler"
	"github.com/nekr0z/muhadi/internal/reconciler/mocks"
	"go.uber.org/mock/gomock"
)

func TestReconciler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	st := mocks.NewMockOrderStorage(ctrl)
	ac := mocks.NewMockAccrual(ctrl)

	ctx, cancel := context.WithCancel(context.Background())

	r := reconciler.New(st, ac)

	gomock.InOrder(
		// normal reconciliation
		st.EXPECT().FirstInQueue(gomock.Any()).Return(newOrder(t, 12345), nil),
		ac.EXPECT().Status(gomock.Any(), 12345).Return(200.0, nil),
		st.EXPECT().UpdateOrder(gomock.Any(), gomock.Cond(func(o *order.Order) bool {
			return o.ID == 12345 && o.Status == order.StatusDone && o.Accrual == 200
		})).Return(nil),

		// empty queue
		st.EXPECT().FirstInQueue(gomock.Any()).Return(nil, fmt.Errorf("nothing here")),

		// invalid order
		st.EXPECT().FirstInQueue(gomock.Any()).Return(newOrder(t, 12344), nil),
		ac.EXPECT().Status(gomock.Any(), 12344).Return(0.0, reconciler.ErrAccrualRejected),
		st.EXPECT().UpdateOrder(gomock.Any(), gomock.Cond(func(o *order.Order) bool {
			return o.ID == 12344 && o.Status == order.StatusInvalid && o.Accrual == 0
		})).Return(nil),

		// not ready
		st.EXPECT().FirstInQueue(gomock.Any()).Return(newOrder(t, 12343), nil),
		ac.EXPECT().Status(gomock.Any(), 12343).Return(0.0, reconciler.ErrAccrualNotReady),
		st.EXPECT().UpdateOrder(gomock.Any(), gomock.Cond(func(o *order.Order) bool {
			return o.ID == 12343 && o.Status == order.StatusProcessing && o.Accrual == 0
		})).Return(nil),

		// backoff
		st.EXPECT().FirstInQueue(gomock.Any()).Return(newOrder(t, 12346), nil),
		ac.EXPECT().Status(gomock.Any(), 12346).Return(0.0, &reconciler.ErrAccrualOverload{60 * time.Second}),
	)

	go r.Run(ctx)

	time.Sleep(2 * time.Second)
	cancel()
}

func newOrder(t *testing.T, id int) *order.Order {
	t.Helper()
	return order.New(id, "test")
}
