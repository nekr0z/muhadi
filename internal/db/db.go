package db

import (
	"context"
	"database/sql"
	"embed"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/nekr0z/muhadi/internal/ctxlog"
	"go.uber.org/zap"
)

//go:embed migrations/*.sql
var fs embed.FS

type DB struct {
	*sql.DB
}

func New(ctx context.Context, dsn string) (*DB, error) {
	if err := applyMigrations(dsn); err != nil {
		ctxlog.Error(ctx, "failed to apply migrations", zap.Error(err))
		return nil, err
	}

	database, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	return &DB{database}, nil
}

func (db *DB) Close() {
	_ = db.DB.Close()
}

func applyMigrations(dsn string) error {
	src, err := iofs.New(fs, "migrations")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithSourceInstance("iofs", src, dsn)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
