package balanceservice

import (
	"context"

	"github.com/nekr0z/muhadi/internal/balance"
)

//go:generate mockgen -destination mocks/storage_mock.go -package mocks . Storage

type Storage interface {
	GetWithdrawals(context.Context, string) ([]balance.Withdrawal, error)
	TotalWithdrawn(context.Context, string) (float64, error)
	SaveWithdrawal(context.Context, *balance.Withdrawal) error
}
