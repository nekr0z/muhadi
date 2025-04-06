package reconciler

import (
	"context"

	"github.com/nekr0z/muhadi/internal/order"
)

//go:generate mockgen -destination mocks/storage_mock.go -package mocks . OrderStorage

type OrderStorage interface {
	FirstInQueue(context.Context) (*order.Order, error)
	UpdateOrder(context.Context, *order.Order) error
}
