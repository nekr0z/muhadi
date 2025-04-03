package logger_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nekr0z/muhadi/internal/api/middleware/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestLog(t *testing.T) {
	respString := "It works!"
	uri := "http://example.com/foo"

	core, observed := observer.New(zap.DebugLevel)
	log := zap.New(core)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, respString)
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", uri, nil)
	rr := httptest.NewRecorder()

	logger.Log(log)(handler).ServeHTTP(rr, req)
	log.Sync()

	logs := observed.All()
	require.Len(t, logs, 1)

	entry := logs[0]
	assert.Contains(t, entry.Message, "request")

	m := entry.ContextMap()

	assert.Equal(t, http.MethodGet, m["method"])
	assert.Equal(t, uri, m["uri"])
	assert.Equal(t, int64(http.StatusOK), m["status"])
	assert.Equal(t, int64(len(respString)), m["size"])
	assert.Contains(t, m, "duration")
}
