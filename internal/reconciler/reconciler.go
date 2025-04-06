package reconciler

import (
	"context"
	"errors"
	"time"

	"github.com/nekr0z/muhadi/internal/ctxlog"
	"github.com/nekr0z/muhadi/internal/order"
	"go.uber.org/zap"
)

var defaultSleep = 1 * time.Second

type Reconciler struct {
	storage    OrderStorage
	accrual    Accrual
	emptySleep time.Duration
}

func New(os OrderStorage, acc Accrual) *Reconciler {
	r := &Reconciler{
		storage:    os,
		accrual:    acc,
		emptySleep: defaultSleep,
	}

	return r
}

func (r *Reconciler) Run(ctx context.Context) {
	ctxlog.Debug(ctx, "Starting order reconciler.")
	for {
		select {
		case <-ctx.Done():
			ctxlog.Debug(ctx, "Stopping order reconciler.")
			return
		default:
			r.reconcile(ctx)
		}
	}
}

func (r *Reconciler) reconcile(ctx context.Context) {
	log := ctxlog.Maybe(ctx)

	for {
		ord, err := r.storage.FirstInQueue(ctx)
		if err != nil {
			log.Debug("Failed to get an order to reconcile.", zap.Error(err))
			time.Sleep(r.emptySleep)
			return
		}

		accrualAmount, err := r.accrual.Status(ctx, ord.ID)
		if err == nil {
			ord.UpdateStatus(order.StatusDone)
			ord.UpdateAccrual(accrualAmount)
		} else if errors.Is(err, ErrAccrualNotReady) {
			ord.UpdateStatus(order.StatusProcessing)
		} else if errors.Is(err, ErrAccrualRejected) {
			ord.UpdateStatus(order.StatusInvalid)
		} else if e, ok := err.(*ErrAccrualOverload); ok {
			log.Info("Accrual service overloaded", zap.Duration("backoff", e.Duration))
			time.Sleep(e.Duration)
			return
		} else {
			log.Warn("Failed to query accrual service", zap.Error(err))
			time.Sleep(r.emptySleep)
			return
		}

		err = r.storage.UpdateOrder(ctx, ord)
		if err != nil {
			log.Warn("Failed to update order status", zap.Int("orderID", ord.ID), zap.Error(err))
		}

		log.Debug("Updated order status", zap.Int("orderID", ord.ID), zap.String("status", ord.Status.String()))
	}
}
