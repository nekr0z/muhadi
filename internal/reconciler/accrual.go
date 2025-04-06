package reconciler

import (
	"context"
	"errors"
	"fmt"
	"time"
)

//go:generate mockgen -destination mocks/accrual_mock.go -package mocks . Accrual

type Accrual interface {
	Status(ctx context.Context, order int) (result float64, err error)
}

type ErrAccrualOverload struct {
	time.Duration
}

func (e *ErrAccrualOverload) Error() string {
	return fmt.Sprintf("requested timeout %s", e.String())
}

var ErrAccrualRejected = errors.New("order permanently rejected")

var ErrAccrualNotReady = errors.New("order processing still in progress")
