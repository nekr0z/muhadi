package ctxlog_test

import (
	"context"
	"testing"

	"github.com/nekr0z/muhadi/internal/ctxlog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestLevels(t *testing.T) {
	core, logs := observer.New(zap.DebugLevel)
	log := zap.New(core)

	ctx := ctxlog.New(context.Background(), log)

	ctxlog.Debug(ctx, "debug message")
	ctxlog.Info(ctx, "info message")
	ctxlog.Warn(ctx, "warning message")
	ctxlog.Error(ctx, "error message")

	entries := logs.All()

	require.Equal(t, len(entries), 4)

	assert.Equal(t, "debug message", entries[0].Message)
	assert.Equal(t, zapcore.DebugLevel, entries[0].Level)

	assert.Equal(t, "info message", entries[1].Message)
	assert.Equal(t, zapcore.InfoLevel, entries[1].Level)

	assert.Equal(t, "warning message", entries[2].Message)
	assert.Equal(t, zapcore.WarnLevel, entries[2].Level)

	assert.Equal(t, "error message", entries[3].Message)
	assert.Equal(t, zapcore.ErrorLevel, entries[3].Level)
}
