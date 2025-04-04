package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nekr0z/muhadi/internal/accrualclient"
	"github.com/nekr0z/muhadi/internal/api"
	balanceservice "github.com/nekr0z/muhadi/internal/balance/service"
	"github.com/nekr0z/muhadi/internal/ctxlog"
	"github.com/nekr0z/muhadi/internal/db"
	orderservice "github.com/nekr0z/muhadi/internal/order/service"
	"github.com/nekr0z/muhadi/internal/reconciler"
	"github.com/nekr0z/muhadi/internal/user"
	"go.uber.org/zap"
)

type App struct {
	db         *db.DB
	server     *http.Server
	reconciler *reconciler.Reconciler
	logger     *zap.Logger
}

func New() *App {
	cfg := newConfig()

	log, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	database, err := db.New(cfg.Database)
	if err != nil {
		panic(err)
	}

	userSvc := user.NewService(database)
	orderSvc := orderservice.New(database)
	balanceSvc := balanceservice.New(database, orderSvc)

	server := &http.Server{
		Addr:    cfg.Listen,
		Handler: api.New(log, userSvc, orderSvc, balanceSvc, balanceSvc),
	}

	accrualClient := accrualclient.New(cfg.Accrual)

	reconciler := reconciler.New(database, accrualClient)

	log.Info("app configured",
		zap.String("listen", cfg.Listen),
		zap.String("database", cfg.Database),
		zap.String("accrual", cfg.Accrual),
	)

	return &App{
		db:         database,
		server:     server,
		reconciler: reconciler,
		logger:     log,
	}
}

func (a *App) Run() error {
	var appError error

	serverChan := make(chan struct{}, 1)
	go func() {
		a.logger.Info("running server")
		if err := a.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			a.logger.Error("HTTP server error", zap.Error(err))
			appError = errors.Join(appError, err)
		}
		a.logger.Info("server stopped")
		serverChan <- struct{}{}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	ctx = ctxlog.New(ctx, a.logger)
	reconcilerChan := make(chan struct{})
	go func() {
		a.logger.Info("running reconciler")
		a.reconciler.Run(ctx)
		a.logger.Info("reconciler stopped")
		close(reconcilerChan)
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigChan:
		a.logger.Info("Signal received, will shutdown")
	case <-serverChan:
		a.logger.Info("Server stopped, will exit")
	}

	cancel()

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownRelease()

	if err := a.server.Shutdown(shutdownCtx); err != nil {
		a.logger.Error("server shutdown error", zap.Error(err))
		appError = errors.Join(appError, err)
	}

	<-reconcilerChan

	if err := a.db.Close(); err != nil {
		a.logger.Error("database close error", zap.Error(err))
		appError = errors.Join(appError, err)
	}

	a.logger.Info("Shutdown complete.")

	return appError
}
