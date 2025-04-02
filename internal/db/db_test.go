package db_test

import (
	"context"
	"testing"
	"time"

	"github.com/nekr0z/muhadi/internal/db"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testDB *db.DB

func TestMain(m *testing.M) {
	ctx := context.Background()

	dbName := "gophermart"
	dbUser := "user"
	dbPassword := "password"

	postgresContainer, err := postgres.Run(ctx,
		"postgres:17",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
			wait.ForListeningPort("5432/tcp")),
	)
	if err != nil {
		panic(err)
	}

	defer func() {
		err := testcontainers.TerminateContainer(postgresContainer)
		if err != nil {
			panic(err)
		}
	}()

	dsn, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		panic(err)
	}

	testDB, err = db.New(ctx, dsn)
	if err != nil {
		panic(err)
	}

	m.Run()
}
