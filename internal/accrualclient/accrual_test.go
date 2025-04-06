package accrualclient_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nekr0z/muhadi/internal/accrualclient"
	"github.com/nekr0z/muhadi/internal/ctxlog"
	"github.com/nekr0z/muhadi/internal/reconciler"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestAccrualClient(t *testing.T) {
	s := httptest.NewServer(testHandler)
	defer s.Close()

	ctx := ctxlog.New(context.Background(), zaptest.NewLogger(t))

	client := accrualclient.New(s.URL)

	t.Run("success", func(t *testing.T) {
		res, err := client.Status(ctx, 12345)
		assert.NoError(t, err)
		assert.Equal(t, 500.0, res)
	})

	t.Run("rejected", func(t *testing.T) {
		_, err := client.Status(ctx, 12346)
		assert.True(t, errors.Is(err, reconciler.ErrAccrualRejected))
	})

	t.Run("processing", func(t *testing.T) {
		_, err := client.Status(ctx, 12347)
		assert.True(t, errors.Is(err, reconciler.ErrAccrualNotReady))
	})

	t.Run("unknown", func(t *testing.T) {
		_, err := client.Status(ctx, 12348)
		assert.True(t, errors.Is(err, reconciler.ErrAccrualNotReady))
	})

	t.Run("overload", func(t *testing.T) {
		_, err := client.Status(ctx, 12349)
		e, ok := err.(*reconciler.ErrAccrualOverload)
		assert.True(t, ok)
		assert.Equal(t, 120.0, e.Duration.Seconds())
	})
}

var testHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/api/orders/12345":
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"order": "12345", "status": "PROCESSED", "accrual": 500}`))
	case "/api/orders/12346":
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"order": "12346", "status": "INVALID"}`))
	case "/api/orders/12347":
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"order": "12347", "status": "PROCESSING"}`))
	case "/api/orders/12348":
		w.WriteHeader(http.StatusNoContent)
	case "/api/orders/12349":
		w.Header().Add("Content-Type", "text/plain")
		w.Header().Add("Retry-After", "120")
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(`No more than 3 requests per minute allowed`))

	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}
})
