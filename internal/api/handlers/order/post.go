package order

import (
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/nekr0z/muhadi/internal/api/handlers"
	"github.com/nekr0z/muhadi/internal/ctxlog"
	"github.com/nekr0z/muhadi/internal/order"
	orderservice "github.com/nekr0z/muhadi/internal/order/service"
	"go.uber.org/zap"
)

func PostOrderHandleFunc(log *zap.Logger, os OrderService, ts handlers.TokenService) func(http.ResponseWriter, *http.Request) {
	log = log.With(zap.String("handler", "user/orders::POST"))

	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		userName, ok := handlers.AuthUser(r, ts)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		log = log.With(zap.String("user", userName))

		b, err := io.ReadAll(r.Body)
		if err != nil {
			log.Error("failed to read request body", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		orderID, err := strconv.Atoi(string(b))
		if err != nil {
			log.Info("failed to parse order ID", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx := ctxlog.New(r.Context(), log)

		err = os.NewOrder(ctx, orderID, userName)

		if errors.Is(err, orderservice.ErrAlreadyExists) {
			log.Info("order already submitted")
			w.WriteHeader(http.StatusOK)
			return
		}

		if errors.Is(err, orderservice.ErrOrderIDAlreadyTaken) {
			log.Info("order submitted by another user")
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}

		if errors.Is(err, order.ErrOrderNumberInvalid) {
			log.Info("order number is invalid")
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		if err != nil {
			log.Error("failed to submit order", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Info("order submitted")
		w.WriteHeader(http.StatusAccepted)
	}
}
