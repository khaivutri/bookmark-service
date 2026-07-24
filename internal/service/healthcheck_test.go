package service

import (
	"context"
	"errors"
	"testing"

	"github.com/khaivutri/bookmark-service/internal/model"
	"github.com/khaivutri/bookmark-service/internal/repository"
	"github.com/stretchr/testify/assert"
)

// stubPinger is a lightweight test double for the unexported
// redisPinger and dbPinger interfaces.
type stubPinger struct {
	err error
}

func (s *stubPinger) Ping(_ context.Context) error {
	return s.err
}

func TestHealthCheck_Check(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           		string
		serviceName    		string
		instanceID     		string
		redisErr       		error
		dbErr          		error
		expectedReport 		*model.HealthReport
		expectErr      		bool
	}{
		{
			name:        		"redis up, db up - healthy report",
			serviceName: 		"bookmark_service",
			instanceID:  		"2947cc38-7c27-4b15-9d2f-50e52e638935",
			redisErr:   		nil,
			dbErr:      		nil,
			expectedReport: 	&model.HealthReport{
				Message:      		"OK",
				ServiceName:  		"bookmark_service",
				InstanceID:  		"2947cc38-7c27-4b15-9d2f-50e52e638935",
				Dependencies:		map[string]string{"redis": "UP", "postgres": "UP"},
			},
			expectErr:			 false,
		},
		{
			name:        		"redis down, db up - degraded report",
			serviceName:		"bookmark_service",
			instanceID:  		"3c983364-29f3-4501-a7cc-5603e16f6827",
			redisErr:    		errors.New("connection refused"),
			dbErr:       		nil,
			expectedReport: 	&model.HealthReport{
				Message:      		"DEGRADED",
				ServiceName:  		"bookmark_service",
				InstanceID:   		"3c983364-29f3-4501-a7cc-5603e16f6827",
				Dependencies: 		map[string]string{"redis": "DOWN", "postgres": "UP"},
			},
			expectErr: true,
		},
		{
			name:        		"redis up, db down - degraded report",
			serviceName:		"bookmark_service",
			instanceID:  		"3c983364-29f3-4501-a7cc-5603e16f6827",
			redisErr:    		nil,
			dbErr:       		errors.New("db connection refused"),
			expectedReport: 	&model.HealthReport{
				Message:      		"DEGRADED",
				ServiceName:  		"bookmark_service",
				InstanceID:   		"3c983364-29f3-4501-a7cc-5603e16f6827",
				Dependencies: 		map[string]string{"redis": "UP", "postgres": "DOWN"},
			},
			expectErr: true,
		},
		{
			name:        		"redis down, db down - degraded report",
			serviceName:		"bookmark_service",
			instanceID:  		"3c983364-29f3-4501-a7cc-5603e16f6827",
			redisErr:    		errors.New("connection refused"),
			dbErr:       		errors.New("db connection refused"),
			expectedReport: 	&model.HealthReport{
				Message:      		"DEGRADED",
				ServiceName:  		"bookmark_service",
				InstanceID:   		"3c983364-29f3-4501-a7cc-5603e16f6827",
				Dependencies: 		map[string]string{"redis": "DOWN", "postgres": "DOWN"},
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rPinger := &stubPinger{err: tc.redisErr}
			dbPinger := &stubPinger{err: tc.dbErr}
			hc := NewHealthCheck(tc.serviceName, tc.instanceID, rPinger, dbPinger)

			report, err := hc.Check(context.Background())

			assert.Equal(t, tc.expectedReport, report)

			if tc.expectErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, repository.ErrDependencyDown)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}