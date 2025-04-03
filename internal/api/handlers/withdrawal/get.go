package withdrawal

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/nekr0z/muhadi/internal/api/handlers"
	"github.com/nekr0z/muhadi/internal/balance"
	"github.com/nekr0z/muhadi/internal/ctxlog"
	"go.uber.org/zap"
)

func GetWithdrawalsHandleFunc(log *zap.Logger, ws WithdrawalService, ts handlers.TokenService) func(http.ResponseWriter, *http.Request) {
	log = log.With(zap.String("handler", "user/withdrawals"))

	return func(w http.ResponseWriter, r *http.Request) {
		userName, ok := handlers.AuthUser(r, ts)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		log = log.With(zap.String("user", userName))
		ctx := ctxlog.New(r.Context(), log)

		withdrawals, err := ws.GetWithdrawals(ctx, userName)
		if err != nil {
			log.Error("failed to get withdrawals", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(withdrawals) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		resp := convert(withdrawals)

		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			log.Error("failed to encode withdrawals", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}

type withdrawalResponse struct {
	OrderID string    `json:"order"`
	Amount  float64   `json:"sum"`
	At      time.Time `json:"processed_at"`
}

func convert(withdrawals []balance.Withdrawal) []withdrawalResponse {
	result := make([]withdrawalResponse, len(withdrawals))
	for i, w := range withdrawals {
		result[i] = withdrawalResponse{
			OrderID: strconv.Itoa(w.OrderID),
			Amount:  w.Amount,
			At:      w.At,
		}
	}
	return result
}
