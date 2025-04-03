package user

import (
	"encoding/json"
	"net/http"

	"github.com/nekr0z/muhadi/internal/ctxlog"
	"go.uber.org/zap"
)

func LoginHandleFunc(log *zap.Logger, us UserService, ts TokenService) func(http.ResponseWriter, *http.Request) {
	log = log.With(zap.String("handler", "user/login"))

	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var req userRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Info("failed to decode request", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx := ctxlog.New(r.Context(), log)

		ok := us.Auth(ctx, req.UserName, req.Password)
		if !ok {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		log.Info("user logged in", zap.String("user", req.UserName))

		authorize(log, w, req.UserName, ts)
	}
}
