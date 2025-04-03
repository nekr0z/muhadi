package db_test

import (
	"context"
	"testing"
	"time"

	"github.com/nekr0z/muhadi/internal/balance"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithdrawals(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	with1 := &balance.Withdrawal{
		OrderID:  23442,
		UserName: testUserName,
		Amount:   100,
		At:       time.Now().AddDate(0, 0, -1),
	}

	with2 := &balance.Withdrawal{
		OrderID:  26332,
		UserName: testUserName,
		Amount:   5000,
		At:       time.Now().AddDate(0, 0, -4),
	}

	t.Run("save withdrawals", func(t *testing.T) {
		err := testDB.SaveWithdrawal(ctx, with2)
		require.NoError(t, err)

		err = testDB.SaveWithdrawal(ctx, with1)
		require.NoError(t, err)
	})

	t.Run("get withdrawals", func(t *testing.T) {
		withdrawals, err := testDB.GetWithdrawals(ctx, testUserName)
		require.NoError(t, err)

		require.Len(t, withdrawals, 2)
		assertSameWithdrawal(t, *with1, withdrawals[0])
		assertSameWithdrawal(t, *with2, withdrawals[1])
	})

	t.Run("total withdrawal", func(t *testing.T) {
		total, err := testDB.TotalWithdrawn(ctx, testUserName)
		assert.NoError(t, err)
		assert.Equal(t, float64(5100), total)
	})

	t.Run("total withdrawal zero", func(t *testing.T) {
		total, err := testDB.TotalWithdrawn(ctx, "other")
		assert.NoError(t, err)
		assert.Equal(t, float64(0), total)
	})
}

func assertSameWithdrawal(t *testing.T, expected, actual balance.Withdrawal) {
	assert.Equal(t, expected.OrderID, actual.OrderID)
	assert.Equal(t, expected.UserName, actual.UserName)
	assert.Equal(t, expected.Amount, actual.Amount)
	assert.WithinDuration(t, expected.At, actual.At, time.Microsecond)
}
