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
		ctx         context.Context
		wantErr     bool
	}{
		{
			name: "database is up -> Ping returns nil",
			setupPinger: func(t *testing.T) *DBPinger {
				return NewDBPinger(sqldb.InitMockDB(t))
			},
		},
		{
			name: "database is closed -> Ping returns error",
			setupPinger: func(t *testing.T) *DBPinger {
				db := sqldb.InitMockDB(t)
				sqlDB, _ := db.DB()
				_ = sqlDB.Close()
				return NewDBPinger(db)
			},
			wantErr: true,
		},
		{
			name: "nil database -> Ping returns error",
			setupPinger: func(t *testing.T) *DBPinger {
				return NewDBPinger(nil)
			},
			wantErr: true,
		},
		{
			name: "context already expired -> Ping returns error",
			setupPinger: func(t *testing.T) *DBPinger {
				return NewDBPinger(sqldb.InitMockDB(t))
			},
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pinger := tc.setupPinger(t)
			ctx := tc.ctx
			if ctx == nil {
				ctx = t.Context()
			}
			err := pinger.Ping(ctx)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

