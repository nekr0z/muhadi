package withdrawal_test

import "net/http"

var (
	testUserName = "test_user"
	cookie       = &http.Cookie{Name: "token", Value: "test_token"}
)

type mockTokenService struct{}

func (t mockTokenService) AuthUser(token string) (string, bool) {
	return testUserName, true
}
