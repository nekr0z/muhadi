package orderservice

import (
	"context"

	"github.com/nekr0z/muhadi/internal/order"
)

//go:generate mockgen -destination mocks/storage_mock.go -package mocks . Storage

type Storage interface {
	CreateOrder(context.Context, *order.Order) error
	GetOrder(ctx context.Context, id int) (*order.Order, error)
	GetOrders(ctx context.Context, userName string) ([]*order.Order, error)
	TotalAccrual(ctx context.Context, userName string) (float64, error)
}
