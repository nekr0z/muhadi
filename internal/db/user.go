package db

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/nekr0z/muhadi/internal/user"
)

var _ user.Storage = &DB{}

const usersTable = "users"

var (
	addUserQuery = fmt.Sprintf("INSERT INTO %s (username, hash) VALUES ($1, $2)", usersTable)
	getHashQuery = fmt.Sprintf("SELECT hash FROM %s WHERE username = $1", usersTable)
)

func (db *DB) NewUser(ctx context.Context, userName string, hash [32]byte) error {
	hashVal := hex.EncodeToString(hash[:])

	_, err := db.ExecContext(ctx, addUserQuery, userName, hashVal)
	if isUniqueViolation(err) {
		return user.ErrAlreadyExists
	}

	return err
}

func (db *DB) GetHash(ctx context.Context, userName string) ([32]byte, error) {
	r := db.QueryRowContext(ctx, getHashQuery, userName)

	var hash string

	err := r.Scan(&hash)
	if err != nil {
		return [32]byte{}, err
	}

	return decode(hash)
}

func decode(hash string) ([32]byte, error) {
	var h [32]byte

	_, err := hex.Decode(h[:], []byte(hash))
	if err != nil {
		return [32]byte{}, err
	}

	return h, nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return false
	}
	return pgErr.Code == pgerrcode.UniqueViolation
}
