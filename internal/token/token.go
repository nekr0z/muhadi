package token

import (
	"crypto/rand"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nekr0z/muhadi/internal/api/handlers"
	"github.com/nekr0z/muhadi/internal/api/handlers/user"
)

var _ user.TokenService = &TokenService{}
var _ handlers.TokenService = &TokenService{}

type TokenService struct {
	secret []byte
}

var tokenExpiration = time.Hour * 4

func New() *TokenService {
	b := make([]byte, 64)
	rand.Read(b)

	return &TokenService{
		secret: b,
	}
}

func (ts *TokenService) NewToken(userName string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   userName,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExpiration)),
	})

	return token.SignedString(ts.secret)
}

func (ts *TokenService) AuthUser(token string) (string, bool) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return ts.secret, nil
	}, jwt.WithValidMethods([]string{"HS256"}))

	if err != nil {
		return "", false
	}

	username, err := t.Claims.GetSubject()

	return username, err == nil
}
