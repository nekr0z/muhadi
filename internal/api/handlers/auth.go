package handlers

import "net/http"

type TokenService interface {
	AuthUser(token string) (string, bool)
}

const TokenCookieName = "token"

func AuthUser(r *http.Request, ts TokenService) (string, bool) {
	c := r.CookiesNamed(TokenCookieName)
	if len(c) == 0 {
		return "", false
	}

	token := c[0].Value

	return ts.AuthUser(token)
}
