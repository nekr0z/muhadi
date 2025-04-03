package db

import (
	"context"
	"fmt"

	"github.com/nekr0z/muhadi/internal/balance"
	balanceservice "github.com/nekr0z/muhadi/internal/balance/service"
	"github.com/nekr0z/muhadi/internal/ctxlog"
	"go.uber.org/zap"
)

var _ balanceservice.Storage = &DB{}

const withdrawalsTable = "withdrawals"

var (
	withdrawalsColOrderID  = "id"
	withdrawalsColUserName = "username"
	withdrawalsColAmount   = "amount"
	withdrawalsColAt       = "at"
)

var (
	withdrawalColsStub = fmt.Sprintf(
		"%s, %s, %s, %s",
		withdrawalsColOrderID,
		withdrawalsColUserName,
		withdrawalsColAmount,
		withdrawalsColAt,
	)
	saveWithdrawalQuery = fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES ($1, $2, $3, $4)",
		withdrawalsTable,
		withdrawalColsStub,
	)
	getWithdrawalsQuery = fmt.Sprintf(
		"SELECT %s FROM %s WHERE %s = $1 ORDER BY %s DESC",
		withdrawalColsStub,
		withdrawalsTable,
		withdrawalsColUserName,
		withdrawalsColAt,
	)
	getTotalWithdrawnQuery = fmt.Sprintf(
		"SELECT SUM(%s) FROM %s WHERE %s = $1",
		withdrawalsColAmount,
		withdrawalsTable,
		withdrawalsColUserName,
	)
)

func (db *DB) SaveWithdrawal(ctx context.Context, withdrawal *balance.Withdrawal) error {
	_, err := db.ExecContext(ctx, saveWithdrawalQuery, withdrawal.OrderID, withdrawal.UserName, withdrawal.Amount, withdrawal.At)
	return err
}

func (db *DB) GetWithdrawals(ctx context.Context, userName string) ([]balance.Withdrawal, error) {
	r, err := db.QueryContext(ctx, getWithdrawalsQuery, userName)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	withdrawals := []balance.Withdrawal{}

	for r.Next() {
		wd := balance.Withdrawal{}
		err := r.Scan(&wd.OrderID, &wd.UserName, &wd.Amount, &wd.At)
		if err != nil {
			ctxlog.Error(ctx, "failed to scan withdrawal", zap.Error(err))
			return withdrawals, err
		}
		withdrawals = append(withdrawals, wd)
	}

	return withdrawals, r.Err()
}

func (db *DB) TotalWithdrawn(ctx context.Context, userName string) (float64, error) {
	r := db.QueryRowContext(ctx, getTotalWithdrawnQuery, userName)

	var value float64

	err := r.Scan(&value)
	return value, err
}
