package withdrawal

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/nekr0z/muhadi/internal/api/handlers"
	"github.com/nekr0z/muhadi/internal/balance"
	"github.com/nekr0z/muhadi/internal/ctxlog"
	"github.com/nekr0z/muhadi/internal/order"
	"go.uber.org/zap"
)

func WithdrawalHandleFunc(log *zap.Logger, ws WithdrawalService, ts handlers.TokenService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(zap.String("handler", "user/balance/withdraw"))

		defer r.Body.Close()

		userName, ok := handlers.AuthUser(r, ts)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		log = log.With(zap.String("user", userName))

		var req withdrawalRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Info("failed to decode request", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		orderID, err := strconv.Atoi(string(req.OrderID))
		if err != nil {
			log.Info("failed to parse order ID", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log = log.With(zap.Int("order_id", orderID))
		ctx := ctxlog.New(r.Context(), log)

		err = ws.Withdraw(ctx, userName, orderID, req.Sum)

		if errors.Is(err, balance.ErrNotEnoughFunds) {
			log.Info("insufficient funds")
			http.Error(w, err.Error(), http.StatusPaymentRequired)
			return
		}

		if errors.Is(err, order.ErrOrderNumberInvalid) {
			log.Info("order number is invalid")
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		if err != nil {
			log.Error("failed to submit withdrawal", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Info("withdrawal submitted")
		w.WriteHeader(http.StatusOK)
	}
}

type withdrawalRequest struct {
	OrderID string  `json:"order"`
	Sum     float64 `json:"sum"`
}
