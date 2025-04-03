package user

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/nekr0z/muhadi/internal/ctxlog"
	"github.com/nekr0z/muhadi/internal/user"
	"go.uber.org/zap"
)

func RegisterHandleFunc(log *zap.Logger, us UserService, ts TokenService) func(http.ResponseWriter, *http.Request) {
	log = log.With(zap.String("handler", "user/register"))

	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var req userRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Info("failed to decode request", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx := ctxlog.New(r.Context(), log)

		err := us.NewUser(ctx, req.UserName, req.Password)
		if errors.Is(err, user.ErrAlreadyExists) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}

		if err != nil {
			log.Error("failed to register user", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Info("user registered", zap.String("user", req.UserName))

		authorize(log, w, req.UserName, ts)
	}
}
