package dbutils

import (
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
)

func TestParseDBErrorPostgresUniqueConstraints(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		constraint string
		expected   error
	}{
		{name: "username migration constraint", constraint: "uni_username", expected: ErrDuplicateUserName},
		{name: "username gorm constraint", constraint: "uni_users_username", expected: ErrDuplicateUserName},
		{name: "email migration constraint", constraint: "uni_email", expected: ErrDuplicateEmail},
		{name: "email gorm constraint", constraint: "uni_users_email", expected: ErrDuplicateEmail},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := &pgconn.PgError{Code: "23505", ConstraintName: tc.constraint}
			assert.ErrorIs(t, ParseDBError(err), tc.expected)
		})
	}
}
