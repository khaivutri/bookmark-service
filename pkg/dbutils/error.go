package dbutils

import (
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type errorFilter = []func(err error) (bool, error)

var (
	ErrDuplicateUserName = errors.New("username already exists")
	ErrDuplicateEmail    = errors.New("email already exists")
	
	ErrDuplicateBookmarkCode = errors.New("bookmark code already exists")

	ErrRecordNotFound    = errors.New("record not found")
)

func filterDuplicateBookmarkCode(err error) (bool, error) {
	// PostgreSQL
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return pgErr.ConstraintName == "uni_bookmarks_code", ErrDuplicateBookmarkCode
	}
	// SQLite: "UNIQUE constraint failed: bookmarks.code"
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "unique constraint failed") &&
		strings.Contains(msg, "bookmarks.code"), ErrDuplicateBookmarkCode
}


func filterDuplicateUserName(err error) (bool, error) {
	// PostgreSQL
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return pgErr.ConstraintName == "uni_users_username" || pgErr.ConstraintName == "uni_users_uusername", ErrDuplicateUserName
	}
	// SQLite: "UNIQUE constraint failed: users.username" or "UNIQUE constraint failed: users.user_name"
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "unique constraint failed") &&
		(strings.Contains(msg, "users.username") || strings.Contains(msg, "users.user_name")), ErrDuplicateUserName
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

func filterRecordNotFound(err error) (bool, error) {
	// GORM: record not found (First/Take/Last failed to find row)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return true, ErrRecordNotFound
	}
	return false, nil
}

var filters errorFilter = []func(error) (bool, error){
	filterDuplicateUserName,
	filterDuplicateEmail,

	filterDuplicateBookmarkCode,

	filterRecordNotFound,
}

// ParseDBError maps raw database/GORM errors to unified business errors.
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
