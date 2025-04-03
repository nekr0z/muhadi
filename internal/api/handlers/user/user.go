package user

import (
	"context"
	"net/http"

	"github.com/nekr0z/muhadi/internal/api/handlers"
	"go.uber.org/zap"
)

type userRequest struct {
	UserName string `json:"login"`
	Password string `json:"password"`
}

type UserService interface {
	NewUser(ctx context.Context, userName, password string) error
	Auth(ctx context.Context, userName, password string) bool
}

type TokenService interface {
	NewToken(userName string) (string, error)
}

func authorize(log *zap.Logger, w http.ResponseWriter, userName string, ts TokenService) {
	token, err := ts.NewToken(userName)
	if err != nil {
		log.Error("failed to create token", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  handlers.TokenCookieName,
		Value: token,
	})

	w.WriteHeader(http.StatusOK)
}
