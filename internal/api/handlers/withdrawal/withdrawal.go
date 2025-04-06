package withdrawal

import (
	"context"

	"github.com/nekr0z/muhadi/internal/balance"
)

//go:generate mockgen -destination mocks/os_mock.go -package mocks . WithdrawalService

type WithdrawalService interface {
	Withdraw(ctx context.Context, userName string, orderID int, amount float64) error
	GetWithdrawals(ctx context.Context, userName string) ([]balance.Withdrawal, error)
}
