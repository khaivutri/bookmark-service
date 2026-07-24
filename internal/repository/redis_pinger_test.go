package repository

import (
	"context"
	"testing"

	redisPkg "github.com/khaivutri/bookmark-service/pkg/redis"
	"github.com/stretchr/testify/assert"
)

func TestRedisPinger_Ping(t *testing.T) {
	tests := []struct {
		name        string
		setupPinger func(t *testing.T) *RedisPinger
		ctx         context.Context
		wantErr     bool
	}{
		{
			name: "redis is up -> Ping returns nil",
			setupPinger: func(t *testing.T) *RedisPinger {
				return NewRedisPinger(redisPkg.InitMockRedis(t))
			},
		},
		{
			name: "redis is down -> Ping returns error",
			setupPinger: func(t *testing.T) *RedisPinger {
				client := redisPkg.InitMockRedis(t)
				_ = client.Close()
				return NewRedisPinger(client)
			},
			wantErr: true,
		},
		{
			name: "context already expired -> Ping returns error",
			setupPinger: func(t *testing.T) *RedisPinger {
				return NewRedisPinger(redisPkg.InitMockRedis(t))
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
			p := tc.setupPinger(t)
			ctx := tc.ctx
			if ctx == nil {
				ctx = t.Context()
			}
			if tc.wantErr {
				assert.Error(t, p.Ping(ctx))
			} else {
				assert.NoError(t, p.Ping(ctx))
			}
		})
	}
}
