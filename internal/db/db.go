package db

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5/stdlib"
)

//go:embed migrations/*.sql
var fs embed.FS

type DB struct {
	*sql.DB
}

func New(dsn string) (*DB, error) {
	if err := applyMigrations(dsn); err != nil {
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	database, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	return &DB{database}, nil
}

func (db *DB) Close() error {
	return db.DB.Close()
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
