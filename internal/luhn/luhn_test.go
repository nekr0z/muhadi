package luhn_test

import (
	"testing"

	"github.com/nekr0z/muhadi/internal/luhn"
)

func TestIsValid(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		if !luhn.IsValid(49927398716) {
			t.Error("valid number is invalid")
		}
	})
	t.Run("invalid", func(t *testing.T) {
		if luhn.IsValid(49927398717) {
			t.Error("invalid number is valid")
		}
	})
}
