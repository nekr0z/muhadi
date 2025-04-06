package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nekr0z/muhadi/internal/api/handlers/balance"
	"github.com/nekr0z/muhadi/internal/api/handlers/order"
	"github.com/nekr0z/muhadi/internal/api/handlers/user"
	"github.com/nekr0z/muhadi/internal/api/handlers/withdrawal"
	"github.com/nekr0z/muhadi/internal/api/middleware/logger"
	"github.com/nekr0z/muhadi/internal/token"
	"go.uber.org/zap"
)

func New(
	log *zap.Logger,
	us user.UserService,
	os order.OrderService,
	bs balance.BalanceService,
	ws withdrawal.WithdrawalService,
) http.Handler {
	r := chi.NewRouter()

	r.Use(logger.Log(log))

	ts := token.New()

	r.Post("/api/user/register", user.RegisterHandleFunc(log, us, ts))
	r.Post("/api/user/login", user.LoginHandleFunc(log, us, ts))

	r.Post("/api/user/orders", order.PostOrderHandleFunc(log, os, ts))
	r.Get("/api/user/orders", order.GetOrdersHandleFunc(log, os, ts))

	r.Get("/api/user/balance", balance.BalanceHandleFunc(log, bs, ts))

	r.Post("/api/user/balance/withdraw", withdrawal.WithdrawalHandleFunc(log, ws, ts))
	r.Get("/api/user/withdrawals", withdrawal.GetWithdrawalsHandleFunc(log, ws, ts))

	return r
}
