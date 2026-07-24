package repository

import (
	"context"
	"testing"

	"github.com/khaivutri/bookmark-service/pkg/sqldb"
	"github.com/stretchr/testify/assert"
)

func TestDBPinger_Ping(t *testing.T) {
	tests := []struct {
		name        string
		setupPinger func(t *testing.T) *DBPinger
		setupCtx    func(t *testing.T) context.Context
		wantErr     bool
	}{
		{
			name: "database is up -> Ping returns nil",
			setupPinger: func(t *testing.T) *DBPinger {
				db := sqldb.InitMockDB(t)
				return NewDBPinger(db)
			},
			setupCtx: func(t *testing.T) context.Context {
				return context.Background()
			},
			wantErr: false,
		},
		{
			name: "database is closed -> Ping returns error",
			setupPinger: func(t *testing.T) *DBPinger {
				db := sqldb.InitMockDB(t)
				sqlDB, err := db.DB()
				if err != nil {
					t.Fatalf("failed to get sqlDB: %v", err)
				}
				_ = sqlDB.Close()
				return NewDBPinger(db)
			},
			setupCtx: func(t *testing.T) context.Context {
				return context.Background()
			},
			wantErr: true,
		},
		{
			name: "nil database -> Ping returns ErrDependencyDown",
			setupPinger: func(t *testing.T) *DBPinger {
				return NewDBPinger(nil)
			},
			setupCtx: func(t *testing.T) context.Context {
				return context.Background()
			},
			wantErr: true,
		},
		{
			name: "context already expired -> Ping returns error",
			setupPinger: func(t *testing.T) *DBPinger {
				db := sqldb.InitMockDB(t)
				return NewDBPinger(db)
			},
			setupCtx: func(t *testing.T) context.Context {
				ctx, cancel := context.WithTimeout(context.Background(), 0) // expires immediately
				t.Cleanup(cancel)
				return ctx
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pinger := tc.setupPinger(t)
			ctx := tc.setupCtx(t)

			err := pinger.Ping(ctx)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
