package orderservice

import (
	"context"
	"errors"

	"github.com/nekr0z/muhadi/internal/ctxlog"
	"github.com/nekr0z/muhadi/internal/luhn"
	"github.com/nekr0z/muhadi/internal/order"
	"go.uber.org/zap"
)

type Service struct {
	storage Storage
}

func New(storage Storage) *Service {
	return &Service{storage: storage}
}

func (s *Service) NewOrder(ctx context.Context, id int, userName string) error {
	if !luhn.IsValid(id) {
		return order.ErrOrderNumberInvalid
	}

	newOrder := order.New(id, userName)

	err := s.storage.CreateOrder(ctx, newOrder)
	if err == nil {
		return nil
	}

	ctxlog.Info(ctx, "order already exists")

	oldOrder, err := s.storage.GetOrder(ctx, id)
	if err != nil {
		ctxlog.Error(ctx, "Failed to get order", zap.Error(err))
		return err
	}

	if oldOrder.UserName != userName {
		return ErrOrderIDAlreadyTaken
	}

	return ErrAlreadyExists
}

func (s *Service) GetOrders(ctx context.Context, userName string) ([]*order.Order, error) {
	return s.storage.GetOrders(ctx, userName)
}

func (s *Service) TotalAccrual(ctx context.Context, userName string) (float64, error) {
	return s.storage.TotalAccrual(ctx, userName)
}

var (
	ErrAlreadyExists       = errors.New("order already exists")
	ErrOrderIDAlreadyTaken = errors.New("order ID already submitted by another user")
)
