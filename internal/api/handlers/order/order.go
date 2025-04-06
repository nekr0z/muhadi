package order

import (
	"context"

	"github.com/nekr0z/muhadi/internal/order"
)

//go:generate mockgen -destination mocks/os_mock.go -package mocks . OrderService

type OrderService interface {
	NewOrder(ctx context.Context, orderID int, userName string) error
	GetOrders(ctx context.Context, userName string) ([]*order.Order, error)
}
