package order

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/nekr0z/muhadi/internal/api/handlers"
	"github.com/nekr0z/muhadi/internal/ctxlog"
	"github.com/nekr0z/muhadi/internal/order"
	"go.uber.org/zap"
)

func GetOrdersHandleFunc(log *zap.Logger, os OrderService, ts handlers.TokenService) func(http.ResponseWriter, *http.Request) {
	log = log.With(zap.String("handler", "user/orders::GET"))

	return func(w http.ResponseWriter, r *http.Request) {
		userName, ok := handlers.AuthUser(r, ts)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		log = log.With(zap.String("user", userName))
		ctx := ctxlog.New(r.Context(), log)

		orders, err := os.GetOrders(ctx, userName)
		if err != nil {
			log.Error("failed to get orders", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(orders) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		resp := convert(orders)

		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			log.Error("failed to encode orders", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}

type orderResponse struct {
	ID         string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    float64   `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
}

func orderToResponse(o order.Order) orderResponse {
	return orderResponse{
		ID:         strconv.Itoa(o.ID),
		Status:     o.Status.String(),
		Accrual:    o.Accrual,
		UploadedAt: o.CreatedAt,
	}
}

func convert(orders []*order.Order) []orderResponse {
	result := make([]orderResponse, len(orders))
	for i, o := range orders {
		result[i] = orderToResponse(*o)
	}
	return result
}
