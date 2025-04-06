package db

import (
	"context"
	"fmt"

	"github.com/nekr0z/muhadi/internal/ctxlog"
	"github.com/nekr0z/muhadi/internal/order"
	orderservice "github.com/nekr0z/muhadi/internal/order/service"
	"github.com/nekr0z/muhadi/internal/reconciler"
	"go.uber.org/zap"
)

var _ orderservice.Storage = &DB{}
var _ reconciler.OrderStorage = &DB{}

const ordersTable = "orders"

var (
	ordersColID         = "id"
	ordersColUserName   = "username"
	ordersColStatus     = "status"
	ordersColAccrual    = "accrual"
	ordersColCreatedAt  = "created_at"
	ordersColLastUpdate = "last_updated_at"
)

var (
	ordersColsStub = fmt.Sprintf(
		"%s, %s, %s, %s, %s, %s",
		ordersColID,
		ordersColUserName,
		ordersColStatus,
		ordersColAccrual,
		ordersColCreatedAt,
		ordersColLastUpdate,
	)
	createOrderQuery = fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6)",
		ordersTable,
		ordersColsStub,
	)
	getOrdersStub = fmt.Sprintf(
		"SELECT %s FROM %s",
		ordersColsStub,
		ordersTable,
	)
	getOrderQuery = fmt.Sprintf(
		"%s WHERE %s = $1",
		getOrdersStub,
		ordersColID,
	)
	getOrdersQuery = fmt.Sprintf(
		"%s WHERE %s = $1 ORDER BY %s",
		getOrdersStub,
		ordersColUserName,
		ordersColCreatedAt,
	)
	getTotalAccrualQuery = fmt.Sprintf(
		"SELECT COALESCE(SUM(%s), 0) FROM %s WHERE %s = $1",
		ordersColAccrual,
		ordersTable,
		ordersColUserName,
	)
	getFirstInQueueQuery = fmt.Sprintf(
		"%s WHERE %s NOT IN (%d, %d) ORDER BY %s LIMIT 1",
		getOrdersStub,
		ordersColStatus,
		order.StatusInvalid,
		order.StatusDone,
		ordersColLastUpdate,
	)
	updateOrderQuery = fmt.Sprintf(
		"UPDATE %s SET %s = $1, %s = $2, %s = $3 WHERE %s = $4",
		ordersTable,
		ordersColStatus,
		ordersColAccrual,
		ordersColLastUpdate,
		ordersColID,
	)
)

func (db *DB) CreateOrder(ctx context.Context, order *order.Order) error {
	_, err := db.ExecContext(ctx, createOrderQuery, order.ID, order.UserName, order.Status, order.Accrual, order.CreatedAt, order.LastUpdate)
	return err
}

func (db *DB) GetOrder(ctx context.Context, id int) (*order.Order, error) {
	r := db.QueryRowContext(ctx, getOrderQuery, id)

	o := order.Order{}

	err := r.Scan(&o.ID, &o.UserName, &o.Status, &o.Accrual, &o.CreatedAt, &o.LastUpdate)
	return &o, err
}

func (db *DB) GetOrders(ctx context.Context, userName string) ([]*order.Order, error) {
	r, err := db.QueryContext(ctx, getOrdersQuery, userName)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	orders := []*order.Order{}

	for r.Next() {
		o := order.Order{}
		err := r.Scan(&o.ID, &o.UserName, &o.Status, &o.Accrual, &o.CreatedAt, &o.LastUpdate)
		if err != nil {
			ctxlog.Error(ctx, "failed to scan order", zap.Error(err))
			return orders, err
		}
		orders = append(orders, &o)
	}

	return orders, r.Err()
}

func (db *DB) TotalAccrual(ctx context.Context, userName string) (float64, error) {
	r := db.QueryRowContext(ctx, getTotalAccrualQuery, userName)

	var value float64

	err := r.Scan(&value)
	return value, err
}

func (db *DB) FirstInQueue(ctx context.Context) (*order.Order, error) {
	r := db.QueryRowContext(ctx, getFirstInQueueQuery)
	o := order.Order{}
	err := r.Scan(&o.ID, &o.UserName, &o.Status, &o.Accrual, &o.CreatedAt, &o.LastUpdate)
	return &o, err
}

func (db *DB) UpdateOrder(ctx context.Context, order *order.Order) error {
	_, err := db.ExecContext(ctx, updateOrderQuery, order.Status, order.Accrual, order.LastUpdate, order.ID)
	return err
}
