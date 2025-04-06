package balance

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/nekr0z/muhadi/internal/api/handlers"
	"github.com/nekr0z/muhadi/internal/ctxlog"
	"go.uber.org/zap"
)

func BalanceHandleFunc(log *zap.Logger, bs BalanceService, ts handlers.TokenService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(zap.String("handler", "user/balance"))

		userName, ok := handlers.AuthUser(r, ts)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		log = log.With(zap.String("user", userName))
		ctx := ctxlog.New(r.Context(), log)

		var resp balanceResponse
		var err error

		resp.Current, resp.Withdrawn, err = bs.CurrentAndWithdrawnBalance(ctx, userName)
		if err != nil {
			log.Error("failed to get orders", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			log.Error("failed to encode orders", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

type balanceResponse struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type BalanceService interface {
	CurrentAndWithdrawnBalance(ctx context.Context, userName string) (float64, float64, error)
}
