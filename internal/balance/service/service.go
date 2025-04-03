package balanceservice

import (
	"context"

	"github.com/nekr0z/muhadi/internal/balance"
	"github.com/nekr0z/muhadi/internal/ctxlog"
	"github.com/nekr0z/muhadi/internal/luhn"
	"github.com/nekr0z/muhadi/internal/order"
	"go.uber.org/zap"
)

type Service struct {
	storage       Storage
	accrualGetter AccrualGetter
}

func New(storage Storage, accrualGetter AccrualGetter) *Service {
	return &Service{
		storage:       storage,
		accrualGetter: accrualGetter,
	}
}

func (s *Service) GetWithdrawals(ctx context.Context, userName string) ([]balance.Withdrawal, error) {
	return s.storage.GetWithdrawals(ctx, userName)
}

func (s *Service) CurrentAndWithdrawnBalance(ctx context.Context, userName string) (float64, float64, error) {
	accrual, err := s.accrualGetter.TotalAccrual(ctx, userName)
	if err != nil {
		return 0, 0, err
	}

	withdrawn, err := s.storage.TotalWithdrawn(ctx, userName)
	if err != nil {
		return 0, 0, err
	}

	return accrual - withdrawn, withdrawn, nil
}

func (s *Service) Withdraw(ctx context.Context, userName string, or int, amount float64) error {
	if !luhn.IsValid(or) {
		return order.ErrOrderNumberInvalid
	}

	have, _, err := s.CurrentAndWithdrawnBalance(ctx, userName)
	if err != nil {
		return err
	}

	if have < amount {
		ctxlog.Debug(ctx, "not enough funds", zap.String("user", userName), zap.Float64("have", have), zap.Float64("want", amount))
		return balance.ErrNotEnoughFunds
	}

	w := balance.NewWithdrawal(userName, or, amount)
	return s.storage.SaveWithdrawal(ctx, w)
}

type AccrualGetter interface {
	TotalAccrual(ctx context.Context, userName string) (float64, error)
}
