package dbutils

import (
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
)

type errorFilter = []func(err error) (bool, error)

var (
	ErrDuplicateUserName = errors.New("username already exists")
	ErrDuplicateEmail    = errors.New("email already exists")
	//ErrRecordNotFound    = errors.New("record not found")
)

func filterDuplicateUserName(err error) (bool, error) {
	// PostgreSQL
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return pgErr.ConstraintName == "uni_users_user_name", ErrDuplicateUserName
	}
	// SQLite: "UNIQUE constraint failed: users.username"
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "unique constraint failed") &&
		strings.Contains(msg, "users.user_name"), ErrDuplicateUserName
}

func filterDuplicateEmail(err error) (bool, error) {
	// PostgreSQL
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return pgErr.ConstraintName == "uni_users_email", ErrDuplicateEmail
	}
	// SQLite: "UNIQUE constraint failed: users.email"
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "unique constraint failed") &&
		strings.Contains(msg, "users.email"), ErrDuplicateEmail
}

var filters errorFilter = []func(error) (bool, error){
	filterDuplicateUserName,
	filterDuplicateEmail,
}

func ParseDBError(err error) error {
	if err == nil {
		return nil
	}

	allFilters := filters

	for _, filter := range allFilters {
		if matched, mappedErr := filter(err); matched {
			return mappedErr
		}
	}

	return err
}
