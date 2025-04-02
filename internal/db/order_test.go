package db_test

import (
	"context"
	"testing"
	"time"

	"github.com/nekr0z/muhadi/internal/order"
	"github.com/stretchr/testify/assert"
)

func TestOrders(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	order1 := &order.Order{
		ID:         42,
		UserName:   testUserName,
		Status:     order.StatusDone,
		Accrual:    100,
		CreatedAt:  time.Now().AddDate(0, 0, -2),
		LastUpdate: time.Now(),
	}

	order2 := &order.Order{
		ID:         26,
		UserName:   testUserName,
		Status:     order.StatusNew,
		Accrual:    50,
		CreatedAt:  time.Now().AddDate(0, 0, -1),
		LastUpdate: time.Now(),
	}

	order3 := &order.Order{
		ID:         26,
		UserName:   testUserName,
		Status:     order.StatusNew,
		Accrual:    500,
		CreatedAt:  time.Now().AddDate(0, 0, -4),
		LastUpdate: time.Now(),
	}

	t.Run("create orders", func(t *testing.T) {
		err := testDB.CreateOrder(ctx, order2)
		assert.NoError(t, err)

		err = testDB.CreateOrder(ctx, order1)
		assert.NoError(t, err)

		err = testDB.CreateOrder(ctx, order3)
		assert.Error(t, err)
	})

	t.Run("get order", func(t *testing.T) {
		order, err := testDB.GetOrder(ctx, order1.ID)
		assert.NoError(t, err)
		assertSameOrder(t, order1, order)

		order, err = testDB.GetOrder(ctx, order2.ID)
		assert.NoError(t, err)
		assertSameOrder(t, order2, order)
	})

	t.Run("get orders", func(t *testing.T) {
		orders, err := testDB.GetOrders(ctx, testUserName)
		assert.NoError(t, err)
		assert.Len(t, orders, 2)

		assertSameOrder(t, order1, orders[0])
		assertSameOrder(t, order2, orders[1])
	})

	t.Run("total accrual", func(t *testing.T) {
		total, err := testDB.TotalAccrual(ctx, testUserName)
		assert.NoError(t, err)
		assert.Equal(t, float64(150), total)
	})

	t.Run("first in queue", func(t *testing.T) {
		order, err := testDB.FirstInQueue(ctx)
		assert.NoError(t, err)
		assertSameOrder(t, order2, order)
	})

	t.Run("update order", func(t *testing.T) {
		order2.UpdateAccrual(2000)
		order2.UpdateStatus(order.StatusProcessing)
		err := testDB.UpdateOrder(ctx, order2)
		assert.NoError(t, err)

		order, err := testDB.GetOrder(ctx, order2.ID)
		assert.NoError(t, err)
		assertSameOrder(t, order2, order)
	})
}

func assertSameOrder(t *testing.T, expected, actual *order.Order) {
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.UserName, actual.UserName)
	assert.Equal(t, expected.Status, actual.Status)
	assert.Equal(t, expected.Accrual, actual.Accrual)
	assert.WithinDuration(t, expected.CreatedAt, actual.CreatedAt, time.Microsecond)
	assert.WithinDuration(t, expected.LastUpdate, actual.LastUpdate, time.Microsecond)
}
